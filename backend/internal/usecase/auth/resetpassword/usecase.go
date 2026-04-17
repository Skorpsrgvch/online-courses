package resetpassword

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Token       string
	NewPassword string
}

type TokenData struct {
	UserID    int
	ExpiresAt time.Time
	Used      bool
}

type TokenChecker interface {
	GetResetToken(ctx context.Context, token string) (*TokenData, error)
	MarkTokenUsed(ctx context.Context, token string) error
}

type PasswordUpdater interface {
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
}

type Usecase struct {
	tokenChecker    TokenChecker
	passwordUpdater PasswordUpdater
}

func NewUsecase(tokenChecker TokenChecker, passwordUpdater PasswordUpdater) (*Usecase, error) {
	if tokenChecker == nil || passwordUpdater == nil {
		return nil, errors.New("tokenChecker and passwordUpdater are required")
	}
	return &Usecase{tokenChecker: tokenChecker, passwordUpdater: passwordUpdater}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	if input.Token == "" || input.NewPassword == "" {
		return domain.ErrInvalidCredentials
	}

	if len(input.NewPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	tokenData, err := u.tokenChecker.GetResetToken(ctx, input.Token)
	if err != nil {
		return domain.ErrInvalidCredentials
	}

	if tokenData.Used {
		return errors.New("reset token has already been used")
	}

	if time.Now().UTC().After(tokenData.ExpiresAt) {
		return errors.New("reset token has expired")
	}

	// Хешируем новый пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Обновляем пароль
	if err := u.passwordUpdater.UpdatePassword(ctx, tokenData.UserID, string(passwordHash)); err != nil {
		return err
	}

	// Помечаем токен как использованный
	return u.tokenChecker.MarkTokenUsed(ctx, input.Token)
}
