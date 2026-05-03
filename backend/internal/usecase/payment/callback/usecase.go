package callback

import (
	"context"
	"errors"
	"log"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

func NewUseCase(paymentRepo PaymentRepository, purchaseCreator PurchaseCreator) *UseCase {
	return &UseCase{
		paymentRepo:     paymentRepo,
		purchaseCreator: purchaseCreator,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	payment, err := u.paymentRepo.GetByID(ctx, input.PaymentID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			log.Printf("[Callback] Payment not found: %s", input.PaymentID)
			return &Output{Processed: false}, domain.ErrPaymentNotFound
		}
		return &Output{Processed: false}, err
	}

	log.Printf("[Callback] Received status '%s' for payment %s (current status: %s)", input.Status, input.PaymentID, payment.Status)

	// Если платеж уже обработан, ничего не делаем (идемпотентность)
	if payment.Status == domain.PaymentStatusSucceeded {
		return &Output{Processed: true, Status: string(payment.Status)}, nil
	}

	switch input.Status {
	case "succeeded", "waiting_for_capture":
		// Обновляем статус
		payment.SetSucceeded()
		if err := u.paymentRepo.Update(ctx, payment); err != nil {
			log.Printf("[Callback] Error updating payment status: %v", err)
			return &Output{Processed: false}, err
		}

		// !!! ВАЖНО: Создаем покупку только после успешного обновления статуса
		log.Printf("[Callback] Creating purchase for user %d, course %d", payment.UserID, payment.CourseID)
		if err := u.purchaseCreator.Create(ctx, payment.UserID, payment.CourseID, payment.PaymentID); err != nil {
			// Если ошибка создания покупки (например, дубликат), логируем, но считаем обработку успешной,
			// так как платеж уже прошел.
			log.Printf("[Callback] Warning: failed to create purchase record (might be duplicate): %v", err)
		}

	case "canceled":
		payment.SetCanceled()
		_ = u.paymentRepo.Update(ctx, payment)

	case "failed":
		payment.SetFailed()
		_ = u.paymentRepo.Update(ctx, payment)
	}

	return &Output{
		PaymentID: payment.PaymentID,
		Status:    string(payment.Status),
		Processed: true,
	}, nil
}
