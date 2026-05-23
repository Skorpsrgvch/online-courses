package forgotpassword

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	Email string
}

type Output struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type UserReader interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error)
}

type CodeStore interface {
	CreateCode(ctx context.Context, userID int, code string, expiresAt time.Time) error
	GetLastCodeTime(ctx context.Context, userID int) (time.Time, error)
}

type EmailSender interface {
	SendResetCode(to, code string) error
}

type Usecase struct {
	userReader  UserReader
	codeStore   CodeStore
	emailSender EmailSender
}

const minInterval = 2 * time.Minute

func NewUsecase(userReader UserReader, codeStore CodeStore, emailSender EmailSender) (*Usecase, error) {
	if userReader == nil || codeStore == nil || emailSender == nil {
		return nil, errors.New("dependencies are required")
	}
	return &Usecase{
		userReader:  userReader,
		codeStore:   codeStore,
		emailSender: emailSender,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	successMsg := "Если email зарегистрирован, на него отправлен код."

	if input.Email == "" {
		return &Output{Message: successMsg, ExpiresIn: 15}, nil
	}

	zap.L().Debug("Password recovery requested", zap.String("email", input.Email))

	user, _, err := u.userReader.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			zap.L().Debug("User not found, returning generic success", zap.String("email", input.Email))
		} else {
			zap.L().Error("DB error on user lookup", zap.Error(err))
		}
		return &Output{Message: successMsg, ExpiresIn: 15}, nil
	}

	// Rate Limit
	lastCodeTime, err := u.codeStore.GetLastCodeTime(ctx, user.ID)
	if err == nil && !lastCodeTime.IsZero() {
		timeSinceLast := time.Since(lastCodeTime)
		if timeSinceLast < minInterval {
			waitTime := int((minInterval - timeSinceLast).Seconds())
			zap.L().Warn("Rate limit hit", zap.Int("user_id", user.ID), zap.Int("wait_seconds", waitTime))
			return &Output{
				Message:   fmt.Sprintf("Код уже был отправлен. Подождите %d сек.", waitTime),
				ExpiresIn: 15,
			}, nil
		}
	}

	code, err := generateCode()
	if err != nil {
		zap.L().Error("Failed to generate code", zap.Error(err))
		return nil, errors.New("ошибка генерации кода")
	}

	expiresAt := time.Now().UTC().Add(15 * time.Minute)

	if err := u.codeStore.CreateCode(ctx, user.ID, code, expiresAt); err != nil {
		zap.L().Error("Failed to save code", zap.Error(err))
		return nil, errors.New("ошибка сохранения кода")
	}

	if err := u.emailSender.SendResetCode(user.Email, code); err != nil {
		zap.L().Error("Failed to send email", zap.Error(err), zap.String("email", user.Email))
		return nil, errors.New("ошибка отправки письма")
	}

	zap.L().Info("Recovery code sent", zap.Int("user_id", user.ID))
	return &Output{
		Message:   "Код отправлен на ваш email.",
		ExpiresIn: 15,
	}, nil
}

func generateCode() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := n.Int64() + 100000
	return fmt.Sprintf("%d", code), nil
}
