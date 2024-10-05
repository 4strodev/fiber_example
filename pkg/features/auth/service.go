package auth

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/4strodev/go_monitoring_example/pkg/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	logger   *slog.Logger
	dbClient *gorm.DB
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type RegisterRequest struct {
	Email    string
	Password string
}

func NewAuthService(logger *slog.Logger, dbClient *gorm.DB) AuthService {
	return AuthService{
		logger,
		dbClient,
	}
}

func (s *AuthService) Register(req RegisterRequest) (err error) {
	var user shared.User

	var total int64
	err = s.dbClient.Model(&shared.User{}).Where("email = ?", req.Email).Limit(1).Count(&total).Error
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	if total > 0 {
		return errors.New("user already exists")
	}

	user.ID = uuid.Must(uuid.NewV7())
	user.Email = req.Email
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	s.dbClient.Save(&user)

	return err
}

func (s *AuthService) Login(req LoginRequest) (res LoginResponse, err error) {
	var user shared.User
	err = s.dbClient.First(&user, "email = ?", req.Email).Error
	if err != nil {
		s.logger.Info("user cannot login")
		return res, errors.New("user not found")
	}
	match := PasswordMatch(user.Password, req.Password)
	if !match {
		return res, errors.New("password doesn't match")
	}
	s.logger.Info("user logged in")

	return res, nil
}

func (s *AuthService) Logout() error {
	s.logger.Info("user logged out")
	return nil
}
