package login

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Email    string
	Password string
}

type Output struct {
	User *domain.User
}

type Authenticator interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error)
}

type Usecase struct {
	authenticator Authenticator
}

func NewUsecase(authenticator Authenticator) (*Usecase, error) {
	if authenticator == nil {
		return nil, errors.New("authenticator is required")
	}
	return &Usecase{authenticator: authenticator}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	if input.Email == "" || input.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	zap.L().Debug("Login attempt", zap.String("email", input.Email))

	user, passwordHash, err := u.authenticator.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			zap.L().Debug("Login failed: user not found", zap.String("email", input.Email))
		} else {
			zap.L().Error("Login DB error", zap.Error(err))
		}
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		zap.L().Warn("Login failed: invalid password", zap.String("email", input.Email))
		return nil, domain.ErrInvalidCredentials
	}

	zap.L().Info("Login successful", zap.Int("user_id", user.ID), zap.String("email", user.Email))
	return &Output{User: user}, nil
}
