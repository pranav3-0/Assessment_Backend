package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	FindUserIdBySub(sub string) (string, error)
	UpdateUserProfile(userID uuid.UUID, data models.UserProfileUpdate) error
	GetUserProfile(userID string) (*models.AssessmentUser, error)
	MapUsersToManager(request models.MapUsersToManagerRequest, creator string) error

	GetAllUsers(limit, offset int, role *string, highestRole, userId string) ([]models.UserFullData, int64, error)
	MigrateUsersToKeycloak(userIds []*string) error
	MigrateRoleToKeycloak(roles []string)
	DeleteUser(userID uuid.UUID) error
}

type UserServiceImpl struct {
	userRepo   repository.UserRepository
	clientRepo repository.ClientRepository
	db         *gorm.DB
}

func NewUserService(userRepo repository.UserRepository, clientRepo repository.ClientRepository, db *gorm.DB) UserService {
	return &UserServiceImpl{userRepo: userRepo, clientRepo: clientRepo, db: db}
}

func (us *UserServiceImpl) FindUserIdBySub(sub string) (string, error) {
	return us.userRepo.FindUserIdBySub(sub)
}

func (us *UserServiceImpl) GetUserProfile(userID string) (*models.AssessmentUser, error) {
	user, err := us.userRepo.FindByUserId(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &user, nil
}

func (us *UserServiceImpl) UpdateUserProfile(userID uuid.UUID, data models.UserProfileUpdate) error {
	err := us.userRepo.UpdateUserProfile(userID, data)
	if err != nil {
		return err
	}
	user, err := us.userRepo.FindByUserId(userID.String())
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	if user.AuthUserID == "" {
		return fmt.Errorf("missing Keycloak user ID for %s", userID)
	}
	if data.UserType != nil {
		if *data.UserType != "kid" && *data.UserType != "candidate" {
			return fmt.Errorf("invalid user_type")
		}
	}
	userPayload := map[string]interface{}{
		"firstName": getIfNotNil(data.FirstName),
		"lastName":  getIfNotNil(data.LastName),
		"email":     getIfNotNil(data.Email),
		"phone":     getIfNotNil(data.Phone),
	}

	ctx := context.Background()
	client, err := us.clientRepo.GetByName(ctx, "dhl")
	if err != nil {
		log.Println("Error getting client")
	}

	keycloakService := NewKeycloakService(client.AuthConfig)
	err = keycloakService.UpdateUserInKeycloak(user.AuthUserID, userPayload, data.Roles)
	if err != nil {
		return fmt.Errorf("keycloak sync failed: %w", err)
	}
	return nil
}

func (us *UserServiceImpl) MapUsersToManager(request models.MapUsersToManagerRequest, creator string) error {
	tx := us.db.Begin()
	for _, uid := range request.UserIDs {
		mapping := &models.UserManagerMapping{
			UserID:    uid,
			ManagerID: request.ManagerID,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: creator,
		}
		cm, err := us.userRepo.AddUserToManagerMapping(tx, *mapping)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Println("Mapping Created: ", cm.ID)
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (us *UserServiceImpl) GetAllUsers(limit, offset int, role *string, highestRole, userId string) ([]models.UserFullData, int64, error) {
	if highestRole == "admin" {
		return us.userRepo.FetchAllUsersWithExtPaginated(offset, limit, role)
	}

	if highestRole == "manager" {
		return us.userRepo.FetchUsersManagedBy(offset, limit, userId, role)
	}

	return nil, 0, errors.New("unauthorized: insufficient role")

}

func (us *UserServiceImpl) MigrateUsersToKeycloak(userIds []*string) error {
	users, err := us.userRepo.FetchAllUsers(userIds)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client, err := us.clientRepo.GetByName(ctx, "dhl")
	keycloakService := NewKeycloakService(client.AuthConfig)
	if err != nil {
		log.Println("Error getting client:", err)
		return errors.New("invalid client")
	}
	var errorList []error

	for _, user := range users {
		password := "password"
		createReq := models.RegisterRequest{
			ClientName: "dhl",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Email:      user.Email,
			Phone:      user.Phone,
			Password:   password,
			Roles:      user.Roles,
		}

		keycloakID, err := keycloakService.CreateUser(ctx, createReq)
		if err != nil {
			log.Println("[Error] keycloakService.CreateUser:", user.Email, " : ", err)
			errorList = append(errorList, err)
			continue
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[Error] GenerateFromPassword:", user.Email, " : ", err)
			errorList = append(errorList, errors.New("failed to hash password"))
			continue
		}
		passStr := string(hashedPassword)
		updateData := models.UserProfileUpdate{
			AuthUserID: &keycloakID,
			Password:   &passStr,
		}

		err = us.userRepo.UpdateUserProfile(user.UserID, updateData)
		if err != nil {
			log.Println("[Error] userRepo Update User:", user.Email, " : ", err)
			errorList = append(errorList, err)
			continue
		}
	}

	finalErr := errors.Join(errorList...)
	return finalErr

}

func (us *UserServiceImpl) MigrateRoleToKeycloak(roles []string) {
	ctx := context.Background()
	client, err := us.clientRepo.GetByName(ctx, "dhl")
	if err != nil {
		log.Println("Error getting client")
	}

	keycloakService := NewKeycloakService(client.AuthConfig)
	keycloakService.MigrateRoles(roles)
}

func getIfNotNil(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (us *UserServiceImpl) DeleteUser(userID uuid.UUID) error {
	return us.userRepo.DeleteUser(userID)
}