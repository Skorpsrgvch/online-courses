package register

import (
	"context"
	"errors"
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

	if input.Email == "" || input.FullName == "" || input.Password == "" {
		zap.L().Warn("Registration failed: missing fields")
		return domain.ErrInvalidCredentials
	}

	if input.Role != "user" && input.Role != "admin" {
		zap.L().Warn("Registration failed: invalid role", zap.String("role", input.Role))
		return errors.New("invalid role")
	}

	user, err := domain.NewUser(input.Email, input.FullName, input.Role)
	if err != nil {
		zap.L().Error("Domain validation failed", zap.Error(err))
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
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
