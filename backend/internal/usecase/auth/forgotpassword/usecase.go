package forgotpassword

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	Email string
}

type Output struct {
	Token     string `json:"token"`
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"` // минуты
}

type UserReader interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error)
}

type TokenCreator interface {
	CreateResetToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
}

type Usecase struct {
	userReader  UserReader
	tokenCreator TokenCreator
}

func NewUsecase(userReader UserReader, tokenCreator TokenCreator) (*Usecase, error) {
	if userReader == nil || tokenCreator == nil {
		return nil, errors.New("userReader and tokenCreator are required")
	}
	return &Usecase{userReader: userReader, tokenCreator: tokenCreator}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	if input.Email == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, _, err := u.userReader.GetUserByEmail(ctx, input.Email)
	if err != nil {
		// Не раскрываем существует ли пользователь — возвращаем успех
		return &Output{
			Message:   "Если email зарегистрирован, на него отправлена ссылка для сброса пароля.",
			ExpiresIn: 30,
		}, nil
	}

	// Генерируем токен
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errors.New("failed to generate reset token")
	}
	token := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().UTC().Add(30 * time.Minute)

	if err := u.tokenCreator.CreateResetToken(ctx, user.ID, token, expiresAt); err != nil {
		return nil, err
	}

	// В реальном приложении здесь отправка email.
	// Для тестирования возвращаем токен.
	return &Output{
		Token:     token,
		Message:   "Токен сброшен. Для теста используйте этот токен в /auth/reset-password.",
		ExpiresIn: 30,
	}, nil
}
