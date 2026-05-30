package confirm

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type PaymentRepository interface {
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}

type PurchaseWriter interface {
	Create(ctx context.Context, userID, courseID int, paymentID string) error
}

type Input struct {
	PaymentID string `json:"payment_id"`
}

type Output struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Success   bool   `json:"success"`
}

type UseCase struct {
	paymentRepo    PaymentRepository
	purchaseWriter PurchaseWriter
}

func NewUseCase(paymentRepo PaymentRepository, purchaseWriter PurchaseWriter) *UseCase {
	if paymentRepo == nil || purchaseWriter == nil {
		zap.L().Fatal("Failed to initialize PaymentConfirmUseCase: dependencies are nil")
	}

	return &UseCase{
		paymentRepo:    paymentRepo,
		purchaseWriter: purchaseWriter,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Info("Payment confirmation requested", zap.String("paymentID", input.PaymentID))

	payment, err := u.paymentRepo.GetByID(ctx, input.PaymentID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			zap.L().Warn("Payment not found for confirmation", zap.String("paymentID", input.PaymentID))
			return &Output{Success: false}, domain.ErrPaymentNotFound
		}

		zap.L().Error("Failed to get payment by ID", zap.String("paymentID", input.PaymentID), zap.Error(err))
		return &Output{Success: false}, err
	}

	if payment.IsExpired() {
		zap.L().Warn("Payment expired", zap.String("paymentID", input.PaymentID))

		payment.SetFailed()
		_ = u.paymentRepo.Update(ctx, payment)

		return &Output{Success: false}, domain.ErrPaymentExpired
	}

	if !payment.CanBeCompleted() {
		zap.L().Warn(
			"Payment cannot be completed due to status",
			zap.String("paymentID", input.PaymentID),
			zap.String("status", string(payment.Status)),
		)

		return &Output{Success: false}, domain.ErrPaymentInvalidStatus
	}

	payment.SetSucceeded()

	if err := u.paymentRepo.Update(ctx, payment); err != nil {
		zap.L().Error("Failed to update payment status to succeeded", zap.String("paymentID", input.PaymentID), zap.Error(err))
		return &Output{Success: false}, err
	}

	if err := u.purchaseWriter.Create(ctx, payment.UserID, payment.CourseID, payment.PaymentID); err != nil {
		zap.L().Error("Failed to create or renew purchase after confirmation", zap.String("paymentID", input.PaymentID), zap.Error(err))
		return &Output{Success: false}, err
	}

	zap.L().Info(
		"Payment confirmed and purchase created or renewed",
		zap.String("paymentID", input.PaymentID),
		zap.Int("userID", payment.UserID),
		zap.Int("courseID", payment.CourseID),
	)

	return &Output{
		PaymentID: payment.PaymentID,
		Status:    string(payment.Status),
		Success:   true,
	}, nil
}
