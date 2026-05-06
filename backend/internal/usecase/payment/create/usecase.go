package create

import (
	"context"
	"fmt"
	"log"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

func NewUseCase(paymentRepo PaymentRepository, courseGetter CourseRepository, purchaseChecker PurchaseChecker, paymentGateway PaymentGateway) *UseCase {
	return &UseCase{
		paymentRepo:     paymentRepo,
		courseGetter:    courseGetter,
		purchaseChecker: purchaseChecker,
		paymentGateway:  paymentGateway, // <-- Инициализация поля
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	// 1. Проверка на повторную покупку
	hasPurchased, err := u.purchaseChecker.HasPurchased(ctx, input.UserID, input.CourseID)
	if err != nil {
		return nil, err
	}
	if hasPurchased {
		return nil, domain.ErrPaymentAlreadyPaid
	}

	// 2. Получаем курс из БД, чтобы взять цену и описание
	course, err := u.courseGetter.GetByID(ctx, input.CourseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// Валидация цены
	if course.Price <= 0 {
		return nil, fmt.Errorf("course price is invalid")
	}

	amount := course.Price
	currency := "RUB" // Хардкод или поле из курса
	description := fmt.Sprintf("Оплата курса '%s'", course.Title)

	confirmation := map[string]interface{}{
		"type":       "redirect",
		"return_url": input.ReturnURL,
	}

	// 4. Вызываем шлюз
	gatewayPayment, err := u.paymentGateway.CreatePayment(amount, currency, description, confirmation)
	if err != nil {
		log.Printf("Gateway error: %v", err)
		return nil, fmt.Errorf("failed to create payment in gateway: %w", err)
	}

	// 5. Сохраняем платеж в БД
	payment := domain.NewPayment(
		gatewayPayment.PaymentID,
		input.UserID,
		input.CourseID,
		amount,
		currency,
		description,
	)

	payment.Status = gatewayPayment.Status
	payment.ConfirmationURL = gatewayPayment.ConfirmationURL

	if err := u.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return &Output{
		PaymentID:       payment.PaymentID,
		ConfirmationURL: payment.ConfirmationURL,
		Status:          string(payment.Status),
	}, nil
}
