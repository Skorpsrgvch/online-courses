package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

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
	GetByUserAndCourse(ctx context.Context, userID, courseID int) (*domain.Purchase, error)
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
	paymentGateway  PaymentGateway
}

func NewUseCase(paymentRepo PaymentRepository, courseGetter CourseRepository, purchaseChecker PurchaseChecker, paymentGateway PaymentGateway) *UseCase {
	if paymentRepo == nil || courseGetter == nil || purchaseChecker == nil || paymentGateway == nil {
		zap.L().Fatal("Failed to initialize PaymentCreateUseCase: dependencies are nil")
	}
	return &UseCase{
		paymentRepo:     paymentRepo,
		courseGetter:    courseGetter,
		purchaseChecker: purchaseChecker,
		paymentGateway:  paymentGateway,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Info("Payment creation started", zap.Int("userID", input.UserID), zap.Int("courseID", input.CourseID))

	existingPurchase, err := u.purchaseChecker.GetByUserAndCourse(ctx, input.UserID, input.CourseID)
	if err != nil && !errors.Is(err, domain.ErrPurchaseNotFound) {
		return nil, err
	}

	if existingPurchase != nil {
		expiresAt := existingPurchase.PurchasedAt.Add(365 * 24 * time.Hour)
		if time.Now().Before(expiresAt) {
			return nil, domain.ErrPaymentAlreadyPaid
		}
	}

	// 2. Получаем курс
	course, err := u.courseGetter.GetByID(ctx, input.CourseID)
	if err != nil {
		zap.L().Error("Failed to get course details", zap.Int("courseID", input.CourseID), zap.Error(err))
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	// Валидация цены
	if course.Price <= 0 {
		zap.L().Warn("Attempted to create payment for free or invalid price course", zap.Int("courseID", input.CourseID), zap.Int("price", course.Price))
		return nil, fmt.Errorf("course price is invalid")
	}

	amount := course.Price
	currency := "RUB"
	description := fmt.Sprintf("Оплата курса '%s'", course.Title)

	confirmation := map[string]interface{}{
		"type":       "redirect",
		"return_url": input.ReturnURL,
	}

	// 3. Вызов шлюза
	zap.L().Debug("Calling payment gateway", zap.Int("amount", amount), zap.String("currency", currency))
	gatewayPayment, err := u.paymentGateway.CreatePayment(amount, currency, description, confirmation)
	if err != nil {
		zap.L().Error("Payment gateway error", zap.Error(err))
		return nil, fmt.Errorf("failed to create payment in gateway: %w", err)
	}

	// 4. Сохранение в БД
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
		zap.L().Error("Failed to save payment to database", zap.String("paymentID", payment.PaymentID), zap.Error(err))
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	zap.L().Info("Payment created successfully", zap.String("paymentID", payment.PaymentID), zap.String("confirmationURL", payment.ConfirmationURL))

	return &Output{
		PaymentID:       payment.PaymentID,
		ConfirmationURL: payment.ConfirmationURL,
		Status:          string(payment.Status),
	}, nil
}
