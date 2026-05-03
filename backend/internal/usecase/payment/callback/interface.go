package callback

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type PaymentRepository interface {
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}

type PurchaseCreator interface {
	Create(ctx context.Context, userID, courseID int, paymentID string) error
}

type Input struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"` // status from YooKassa
}

type Output struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Processed bool   `json:"processed"`
}

type UseCase struct {
	paymentRepo     PaymentRepository
	purchaseCreator PurchaseCreator
}
