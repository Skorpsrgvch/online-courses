package create

import (
	"context"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

// PaymentGateway определяет контракт для взаимодействия с платежной системой (ЮKassa)
type PaymentGateway interface {
	CreatePayment(amount int, currency string, description string, confirmation map[string]interface{}) (*domain.Payment, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	UpdateConfirmationURL(ctx context.Context, paymentID string, url string) error
	SetExpiresAt(ctx context.Context, paymentID string, expiresAt time.Time) error
	Update(ctx context.Context, payment *domain.Payment) error
}

type CourseRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}

type PurchaseChecker interface {
	HasPurchased(ctx context.Context, userID, courseID int) (bool, error)
}

type Input struct {
	UserID    int    `json:"user_id"`
	CourseID  int    `json:"course_id"`
	ReturnURL string `json:"return_url"`
}

type Output struct {
	PaymentID       string  `json:"payment_id"`
	ConfirmationURL string  `json:"confirmation_url"`
	Status          string  `json:"status"`
	ExpiresAt       *string `json:"expires_at,omitempty"`
}

type UseCase struct {
	paymentRepo     PaymentRepository
	courseGetter    CourseRepository
	purchaseChecker PurchaseChecker
	paymentGateway  PaymentGateway // <-- Добавлено это поле
}
