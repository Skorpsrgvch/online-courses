package callback

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type PaymentRepository interface {
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}

type PaymentStatusVerifier interface {
	GetPaymentStatus(paymentID string) (status string, paid bool, err error)
}

type PurchaseRepository interface {
	GetByUserAndCourse(ctx context.Context, userID, courseID int) (*domain.Purchase, error)
	Create(ctx context.Context, userID, courseID int, paymentID string) error
	UpdatePurchaseDate(ctx context.Context, userID, courseID int, newDate time.Time) error
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
	paymentRepo  PaymentRepository
	purchaseRepo PurchaseRepository
	verifier     PaymentStatusVerifier
}

func NewUseCase(paymentRepo PaymentRepository, purchaseRepo PurchaseRepository, verifier PaymentStatusVerifier) *UseCase {
	if paymentRepo == nil || purchaseRepo == nil || verifier == nil {
		zap.L().Fatal("Failed to initialize PaymentCallbackUseCase: dependencies are nil")
	}
	return &UseCase{
		paymentRepo:  paymentRepo,
		purchaseRepo: purchaseRepo,
		verifier:     verifier,
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

	if payment.Status == domain.PaymentStatusSucceeded {
		zap.L().Info("Payment already succeeded, skipping processing", zap.String("paymentID", input.PaymentID))
		return &Output{PaymentID: payment.PaymentID, Processed: true, Status: string(payment.Status)}, nil
	}

	status, paid, err := u.verifier.GetPaymentStatus(input.PaymentID)
	if err != nil {
		zap.L().Error("Failed to verify payment status in YooKassa", zap.String("paymentID", input.PaymentID), zap.Error(err))
		return &Output{Processed: false}, err
	}

	if paid && status == "waiting_for_capture" {
		status = "succeeded"
	}

	if input.Status != "" && input.Status != status {
		zap.L().Warn("Webhook status differs from verified YooKassa status",
			zap.String("paymentID", input.PaymentID),
			zap.String("webhookStatus", input.Status),
			zap.String("verifiedStatus", status))
	}

	switch status {
	case "succeeded", "waiting_for_capture":
		payment.SetSucceeded()

		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			zap.L().Error("Failed to update payment status to succeeded", zap.String("paymentID", input.PaymentID), zap.Error(err))
			return &Output{Processed: false}, err
		}

		existingPurchase, err := u.purchaseRepo.GetByUserAndCourse(ctx, payment.UserID, payment.CourseID)
		if err != nil && !errors.Is(err, domain.ErrPurchaseNotFound) {
			zap.L().Error("Failed to check existing purchase", zap.Error(err))
		}

		now := time.Now()
		if existingPurchase != nil {
			if err := u.purchaseRepo.UpdatePurchaseDate(ctx, payment.UserID, payment.CourseID, now); err != nil {
				zap.L().Error("Failed to update purchase date", zap.Error(err))
				return &Output{Processed: false}, err
			}
		} else {
			if err := u.purchaseRepo.Create(ctx, payment.UserID, payment.CourseID, payment.PaymentID); err != nil {
				zap.L().Error("Failed to create purchase", zap.Error(err))
				return &Output{Processed: false}, err
			}
		}

	case "canceled":
		payment.SetCanceled()
		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			return &Output{Processed: false}, err
		}

	case "failed":
		payment.SetFailed()
		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			return &Output{Processed: false}, err
		}

	default:
		zap.L().Warn("Unknown verified payment status", zap.String("status", status), zap.String("paymentID", input.PaymentID))
	}

	return &Output{
		PaymentID: payment.PaymentID,
		Status:    string(payment.Status),
		Processed: true,
	}, nil
}
