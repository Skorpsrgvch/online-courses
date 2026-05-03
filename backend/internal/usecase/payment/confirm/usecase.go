package confirm

import (
	"context"
	"errors"

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
			return &Output{Success: false}, domain.ErrPaymentNotFound
		}
		return &Output{Success: false}, err
	}

	// Проверка на истечение
	if payment.IsExpired() {
		payment.SetFailed()
		_ = u.paymentRepo.Update(ctx, payment)
		return &Output{Success: false}, domain.ErrPaymentExpired
	}

	// Проверка статуса
	if !payment.CanBeCompleted() {
		return &Output{Success: false}, domain.ErrPaymentInvalidStatus
	}

	// Успешное подтверждение
	payment.SetSucceeded()
	if err := u.paymentRepo.Update(ctx, payment); err != nil {
		return &Output{Success: false}, err
	}

	// Создаем запись о покупке
	if err := u.purchaseCreator.Create(ctx, payment.UserID, payment.CourseID, payment.PaymentID); err != nil {
		return &Output{Success: false}, err
	}

	return &Output{
		PaymentID: payment.PaymentID,
		Status:    string(payment.Status),
		Success:   true,
	}, nil
}
