package forgotpassword

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
	GetLastCodeTime(ctx context.Context, userID int) (time.Time, error) // Новый метод
}

type EmailSender interface {
	SendResetCode(to, code string) error
}

type Usecase struct {
	userReader  UserReader
	codeStore   CodeStore
	emailSender EmailSender
}

// Минимальный интервал между запросами кода (в минутах)
const minInterval = 2 * time.Minute

func NewUsecase(userReader UserReader, codeStore CodeStore, emailSender EmailSender) (*Usecase, error) {
	if userReader == nil || codeStore == nil || emailSender == nil {
		return nil, errors.New("dependencies are required")
	}
	return &Usecase{userReader: userReader, codeStore: codeStore, emailSender: emailSender}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	if input.Email == "" {
		// Возвращаем успех, чтобы не светить наличием полей
		return &Output{Message: "Если email зарегистрирован, на него отправлен код.", ExpiresIn: 15}, nil
	}

	user, _, err := u.userReader.GetUserByEmail(ctx, input.Email)
	if err != nil {
		// Если пользователя нет — все равно говорим "успех"
		return &Output{Message: "Если email зарегистрирован, на него отправлен код.", ExpiresIn: 15}, nil
	}

	// === ПРОВЕРКА RATE LIMIT (Бизнес-уровень) ===
	lastCodeTime, err := u.codeStore.GetLastCodeTime(ctx, user.ID)
	if err == nil && !lastCodeTime.IsZero() {
		timeSinceLast := time.Since(lastCodeTime)
		if timeSinceLast < minInterval {
			waitTime := int((minInterval - timeSinceLast).Seconds())
			log.Printf("[RateLimit] User %d tried to request code too soon. Wait %d seconds.", user.ID, waitTime)
			// Возвращаем успех, но код не шлем и не создаем
			return &Output{
				Message:   fmt.Sprintf("Код уже был отправлен. Подождите %d сек.", waitTime),
				ExpiresIn: 15,
			}, nil
		}
	}

	// Генерируем код
	code, err := generateCode()
	if err != nil {
		log.Printf("[ForgotPassword] Failed to generate code: %v", err)
		return nil, errors.New("failed to generate code")
	}

	expiresAt := time.Now().UTC().Add(15 * time.Minute)

	// Сохраняем код (тут же инвалидируются старые)
	if err := u.codeStore.CreateCode(ctx, user.ID, code, expiresAt); err != nil {
		log.Printf("[ForgotPassword] Failed to save code: %v", err)
		return nil, errors.New("failed to save code")
	}

	// Отправляем email
	if err := u.emailSender.SendResetCode(user.Email, code); err != nil {
		log.Printf("[ForgotPassword] Failed to send email: %v", err)
		// В продакшене лучше вернуть ошибку или залогировать, но пользователю сказать успех,
		// чтобы не проверять существование почты через ошибки SMTP.
		// Но для отладки пока вернем ошибку, если SMTP упал.
		return nil, errors.New("failed to send email")
	}

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
