package list

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type PaymentRepository interface {
	ListByUser(ctx context.Context, userID int) ([]*domain.Payment, error)
}

type Input struct {
	UserID int `json:"user_id"`
}

type PaymentItem struct {
	PaymentID       string  `json:"payment_id"`
	CourseID        int     `json:"course_id"`
	Amount          int     `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	PaidAt          *string `json:"paid_at,omitempty"`
	ConfirmationURL string  `json:"confirmation_url,omitempty"`
}

type Output struct {
	Payments []PaymentItem `json:"payments"`
}

type UseCase struct {
	paymentRepo PaymentRepository
}
