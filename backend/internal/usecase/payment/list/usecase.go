package list

import (
	"context"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
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

func NewUseCase(paymentRepo PaymentRepository) *UseCase {
	if paymentRepo == nil {
		zap.L().Fatal("Failed to initialize PaymentListUseCase: repository is nil")
	}
	return &UseCase{
		paymentRepo: paymentRepo,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("Listing payments for user", zap.Int("userID", input.UserID))

	payments, err := u.paymentRepo.ListByUser(ctx, input.UserID)
	if err != nil {
		zap.L().Error("Failed to retrieve payments from repository", zap.Int("userID", input.UserID), zap.Error(err))
		return nil, err
	}

	result := make([]PaymentItem, 0, len(payments))
	for _, p := range payments {
		var paidAtStr *string
		if p.PaidAt != nil {
			s := p.PaidAt.Format(time.RFC3339)
			paidAtStr = &s
		}

		result = append(result, PaymentItem{
			PaymentID:       p.PaymentID,
			CourseID:        p.CourseID,
			Amount:          p.Amount,
			Currency:        p.Currency,
			Status:          string(p.Status),
			CreatedAt:       p.CreatedAt.Format(time.RFC3339),
			PaidAt:          paidAtStr,
			ConfirmationURL: p.ConfirmationURL,
		})
	}

	zap.L().Info("Payments listed successfully", zap.Int("userID", input.UserID), zap.Int("count", len(result)))

	return &Output{
		Payments: result,
	}, nil
}
