package resetpassword

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Code        string
	NewPassword string
}

type TokenData struct {
	UserID    int
	ExpiresAt time.Time
	Used      bool
}

type CodeChecker interface {
	GetByCode(ctx context.Context, code string) (*TokenData, error)
	MarkTokenUsed(ctx context.Context, code string) error
}

type PasswordUpdater interface {
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
}

type Usecase struct {
	codeChecker     CodeChecker
	passwordUpdater PasswordUpdater
}

func NewUsecase(codeChecker CodeChecker, passwordUpdater PasswordUpdater) (*Usecase, error) {
	if codeChecker == nil || passwordUpdater == nil {
		return nil, errors.New("dependencies are required")
	}
	return &Usecase{
		codeChecker:     codeChecker,
		passwordUpdater: passwordUpdater,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("Password reset attempt", zap.Int("code_len", len(input.Code)))

	if input.Code == "" || input.NewPassword == "" {
		return domain.ErrInvalidCredentials
	}

	if len(input.NewPassword) < 6 {
		zap.L().Warn("Password too short")
		return errors.New("пароль должен быть не менее 6 символов")
	}

	tokenData, err := u.codeChecker.GetByCode(ctx, input.Code)
	if err != nil {
		zap.L().Warn("Invalid reset code", zap.Error(err))
		return domain.ErrInvalidCredentials
	}

	if tokenData.Used {
		zap.L().Warn("Reset code already used", zap.Int("user_id", tokenData.UserID))
		return errors.New("код уже был использован")
	}

	if time.Now().UTC().After(tokenData.ExpiresAt) {
		zap.L().Warn("Reset code expired", zap.Int("user_id", tokenData.UserID))
		return errors.New("срок действия кода истек")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Password hashing failed", zap.Error(err))
		return errors.New("ошибка хеширования пароля")
	}

	if err := u.passwordUpdater.UpdatePassword(ctx, tokenData.UserID, string(passwordHash)); err != nil {
		zap.L().Error("Password update failed", zap.Error(err), zap.Int("user_id", tokenData.UserID))
		return err
	}

	if err := u.codeChecker.MarkTokenUsed(ctx, input.Code); err != nil {
		zap.L().Error("Failed to mark code as used", zap.Error(err))
		return err
	}

	zap.L().Info("Password reset successful", zap.Int("user_id", tokenData.UserID))
	return nil
}
