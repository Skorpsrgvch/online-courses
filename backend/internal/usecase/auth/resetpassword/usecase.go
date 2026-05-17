package resetpassword

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Code        string // Теперь это 6-значный код
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
	return &Usecase{codeChecker: codeChecker, passwordUpdater: passwordUpdater}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	if input.Code == "" || input.NewPassword == "" {
		return domain.ErrInvalidCredentials
	}

	if len(input.NewPassword) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}

	// Проверяем код
	tokenData, err := u.codeChecker.GetByCode(ctx, input.Code)
	if err != nil {
		return domain.ErrInvalidCredentials
	}

	if tokenData.Used {
		return errors.New("код уже был использован")
	}

	if time.Now().UTC().After(tokenData.ExpiresAt) {
		return errors.New("срок действия кода истек")
	}

	// Хешируем новый пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	// Обновляем пароль
	if err := u.passwordUpdater.UpdatePassword(ctx, tokenData.UserID, string(passwordHash)); err != nil {
		return err
	}

	// Помечаем код как использованный
	return u.codeChecker.MarkTokenUsed(ctx, input.Code)
}
