package register

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Email    string
	FullName string
	Password string
	Role     string
}

type UserCreator interface {
	CreateUser(ctx context.Context, user *domain.User, passwordHash string) error
}

type Usecase struct {
	userCreator UserCreator
}

func NewUsecase(userCreator UserCreator) (*Usecase, error) {
	if userCreator == nil {
		return nil, errors.New("userCreator is required")
	}
	return &Usecase{userCreator: userCreator}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Info("Registration attempt", zap.String("email", input.Email), zap.String("role", input.Role))

	email := strings.ToLower(strings.TrimSpace(input.Email))
	fullName := strings.TrimSpace(input.FullName)
	password := strings.TrimSpace(input.Password)

	if email == "" || fullName == "" || password == "" {
		zap.L().Warn("Registration failed: missing fields")
		return domain.ErrInvalidCredentials
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email || strings.ContainsAny(email, " \t\r\n") {
		zap.L().Warn("Registration failed: invalid email", zap.String("email", input.Email))
		return errors.New("укажите корректный email")
	}

	if len(password) < 6 {
		zap.L().Warn("Registration failed: weak password")
		return errors.New("пароль должен быть не короче 6 символов")
	}

	if input.Role != "user" && input.Role != "admin" {
		zap.L().Warn("Registration failed: invalid role", zap.String("role", input.Role))
		return errors.New("invalid role")
	}

	user, err := domain.NewUser(email, fullName, input.Role)
	if err != nil {
		zap.L().Error("Domain validation failed", zap.Error(err))
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Password hashing failed", zap.Error(err))
		return errors.New("failed to hash password")
	}

	if err := u.userCreator.CreateUser(ctx, user, string(passwordHash)); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			zap.L().Info("Registration failed: user exists", zap.String("email", input.Email))
			return err
		}
		zap.L().Error("Registration DB error", zap.Error(err))
		return errors.New("failed to register user")
	}

	zap.L().Info("User registered successfully", zap.Int("user_id", user.ID), zap.String("email", user.Email))
	return nil
}
