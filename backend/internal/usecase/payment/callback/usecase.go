package callback

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

type PurchaseCreator interface {
	Create(ctx context.Context, userID, courseID int, paymentID string) error
}

type Input struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
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

func NewUseCase(paymentRepo PaymentRepository, purchaseCreator PurchaseCreator) *UseCase {
	if paymentRepo == nil || purchaseCreator == nil {
		zap.L().Fatal("Failed to initialize PaymentCallbackUseCase: dependencies are nil")
	}
	return &UseCase{
		paymentRepo:     paymentRepo,
		purchaseCreator: purchaseCreator,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Info("Payment callback received", zap.String("paymentID", input.PaymentID), zap.String("status", input.Status))

	payment, err := u.paymentRepo.GetByID(ctx, input.PaymentID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			zap.L().Warn("Payment not found in callback", zap.String("paymentID", input.PaymentID))
			return &Output{Processed: false}, domain.ErrPaymentNotFound
		}
		zap.L().Error("Failed to get payment by ID", zap.String("paymentID", input.PaymentID), zap.Error(err))
		return &Output{Processed: false}, err
	}

	zap.L().Debug("Current payment status", zap.String("currentStatus", string(payment.Status)))

	// Идемпотентность: если уже успешно обработан
	if payment.Status == domain.PaymentStatusSucceeded {
		zap.L().Info("Payment already succeeded, skipping processing", zap.String("paymentID", input.PaymentID))
		return &Output{Processed: true, Status: string(payment.Status)}, nil
	}

	switch input.Status {
	case "succeeded", "waiting_for_capture":
		zap.L().Info("Processing successful payment", zap.String("paymentID", input.PaymentID))
		payment.SetSucceeded()

		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			zap.L().Error("Failed to update payment status to succeeded", zap.String("paymentID", input.PaymentID), zap.Error(err))
			return &Output{Processed: false}, err
		}

		// Создаем покупку
		if err := u.purchaseCreator.Create(ctx, payment.UserID, payment.CourseID, payment.PaymentID); err != nil {
			// Логируем ошибку, но не прерываем процесс, т.к. платеж уже успешен
			zap.L().Warn("Failed to create purchase record (possible duplicate)",
				zap.String("paymentID", input.PaymentID),
				zap.Int("userID", payment.UserID),
				zap.Int("courseID", payment.CourseID),
				zap.Error(err))
		} else {
			zap.L().Info("Purchase created successfully", zap.Int("userID", payment.UserID), zap.Int("courseID", payment.CourseID))
		}

	case "canceled":
		zap.L().Info("Processing canceled payment", zap.String("paymentID", input.PaymentID))
		payment.SetCanceled()
		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			zap.L().Error("Failed to update payment status to canceled", zap.String("paymentID", input.PaymentID), zap.Error(err))
			return &Output{Processed: false}, err
		}

	case "failed":
		zap.L().Info("Processing failed payment", zap.String("paymentID", input.PaymentID))
		payment.SetFailed()
		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			zap.L().Error("Failed to update payment status to failed", zap.String("paymentID", input.PaymentID), zap.Error(err))
			return &Output{Processed: false}, err
		}

	default:
		zap.L().Warn("Unknown payment status received", zap.String("status", input.Status), zap.String("paymentID", input.PaymentID))
	}

	return &Output{
		PaymentID: payment.PaymentID,
		Status:    string(payment.Status),
		Processed: true,
	}, nil
}
