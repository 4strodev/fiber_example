package auth

import (
	"errors"
	"fmt"

	"github.com/4strodev/go_monitoring_example/pkg/shared"
	"github.com/4strodev/go_monitoring_example/pkg/shared/events"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	dbClient *gorm.DB
	eventBus events.EventBus
}

func NewAuthService(dbClient *gorm.DB, eventBus events.EventBus) AuthService {
	return AuthService{
		dbClient,
		eventBus,
	}
}

func (s *AuthService) Register(req RegisterRequest) (err error) {
	var user shared.UserSchema

	var total int64
	err = s.dbClient.Model(&shared.UserSchema{}).Where("email = ?", req.Email).Limit(1).Count(&total).Error
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

	err = s.dbClient.Save(&user).Error
	if err != nil {
		return err
	}
	event := NewUserRegisteredEvent(user.ID.String())
	return s.eventBus.Emit(event)
}

func (s *AuthService) Login(req LoginRequest) (res LoginResponse, err error) {
	var user shared.UserSchema
	err = s.dbClient.First(&user, "email = ?", req.Email).Error
	if err != nil {
		return res, errors.New("user not found")
	}
	match := PasswordMatch(user.Password, req.Password)
	if !match {
		return res, errors.New("password doesn't match")
	}
	err = s.eventBus.Emit(NewUserLoggedInEvent(user.ID.String()))
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *AuthService) Logout() error {
	return nil
}
