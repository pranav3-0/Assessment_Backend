package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"dhl/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req models.RegisterRequest) (*models.AssessmentUser, error)
	LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
	SendOTP(email string) error
	VerifyOTP(email string, otp string) (string, string, error)
}

type AuthServiceImpl struct {
	userRepo            repository.UserRepository
	clientRepo          repository.ClientRepository
	authRepo            repository.AuthRepository
	notificationService NotificationService
	db                  *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	clientRepo repository.ClientRepository,
	authRepo repository.AuthRepository,
	notificationService NotificationService,
	db *gorm.DB,
) AuthService {
	return &AuthServiceImpl{
		userRepo:            userRepo,
		clientRepo:          clientRepo,
		authRepo:            authRepo,
		notificationService: notificationService,
		db:                  db,
	}
}

func (s *AuthServiceImpl) RegisterUser(ctx context.Context, req models.RegisterRequest) (*models.AssessmentUser, error) {
	client, err := s.clientRepo.GetByName(ctx, req.ClientName)
	if err != nil {
		return nil, errors.New("invalid client")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.AssessmentUser{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Username:  req.Email,
		UserType:  req.UserType,
		Password:  string(hashedPassword),
		IsActive:  true,
	}

	keycloakService := NewKeycloakService(client.AuthConfig)
	keycloakID, err := keycloakService.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("keycloak user creation failed: %v", err)
	}
	user.AuthUserID = keycloakID
	user.Password = string(hashedPassword)

	notifyId, err := s.notificationService.RegisterUserInNotify(nil, &user.Phone, user.Email)
	if err == nil {
		user.NotifyId = notifyId.String()
	}

	tx := s.db.Begin()
	if err := s.userRepo.CreateUser(tx, ctx, user, req.Roles); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	if user.UserID == uuid.Nil {
		return nil, errors.New("invalid username or password")
	}
	client, err := s.clientRepo.GetByName(ctx, req.ClientID)
	if err != nil {
		return nil, errors.New("invalid client")
	}

	var response models.LoginResponse
	switch client.AuthType {
	case "keycloak":
		keycloakService := NewKeycloakService(client.AuthConfig)
		accessToken, err := keycloakService.LoginUser(ctx, req)
		if err != nil {
			return nil, err
		}
		response.Token = accessToken.AccessToken
		response.RefreshToken = accessToken.RefreshToken

	case "jwt":

	default:
		return nil, errors.New("unsupported authentication type")
	}

	response.User = &models.UserWithRoles{
		UserID:     user.UserID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Phone:      user.Phone,
		Username:   user.Username,
		AuthUserID: user.AuthUserID,
		Roles:      user.Roles,
		UserType:   user.UserType,
	}

	return &response, nil
}

func (s *AuthServiceImpl) SendOTP(email string) error {
	otp := utils.GenerateOTP()

	err := s.authRepo.SaveOTP(email, otp)
	if err != nil {
		return err
	}

	return utils.SendEmailOTP(email, otp)
}

func (s *AuthServiceImpl) VerifyOTP(email string, otp string) (string, string, error) {
	valid, err := s.authRepo.ValidateOTP(email, otp)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", errors.New("invalid otp")
	}

	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		newUserID := uuid.New()
		err = s.userRepo.CreatePublicUser(newUserID.String(), email)
		if err != nil {
			return "", "", err
		}
		user = &models.AssessmentUser{
			UserID: newUserID,
			Email:  email,
		}
	}

	token, err := GetFreshPublicAssessmentToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to get keycloak token: %w", err)
	}

	return token, user.UserID.String(), nil
}