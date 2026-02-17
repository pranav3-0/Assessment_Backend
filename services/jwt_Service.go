package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	userRepo repository.UserRepository
}

func NewJWTService(userRepo repository.UserRepository) *JWTService {
	return &JWTService{userRepo: userRepo}
}

func (s *JWTService) Authenticate(ctx context.Context, creds map[string]string) (*models.AssessmentUser, error) {
	// user, err := s.userRepo.FindByUsername(ctx, creds["username"])
	// if err != nil {
	// 	return nil, errors.New("invalid username or password")
	// }

	user := models.AssessmentUser{}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds["password"])) != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	// token, err := utils.GenerateJWTToken(user.UserID.String())
	// if err != nil {
	// 	return nil, err
	// }

	user.AuthUserID = "token"
	return &user, nil
}
