package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type fakePaymentRepo struct {
	created *domain.Payment
	err     error
}

func (f *fakePaymentRepo) Create(_ context.Context, payment *domain.Payment) error {
	if f.err != nil {
		return f.err
	}
	f.created = payment
	return nil
}

func (f *fakePaymentRepo) GetByID(context.Context, string) (*domain.Payment, error) { return nil, nil }
func (f *fakePaymentRepo) UpdateConfirmationURL(context.Context, string, string) error { return nil }
func (f *fakePaymentRepo) SetExpiresAt(context.Context, string, time.Time) error { return nil }
func (f *fakePaymentRepo) Update(context.Context, *domain.Payment) error { return nil }

type fakeCourseRepo struct {
	course *domain.Course
	err    error
}

func (f *fakeCourseRepo) GetByID(_ context.Context, _ int) (*domain.Course, error) {
	return f.course, f.err
}

type fakePurchaseChecker struct {
	purchase *domain.Purchase
	err      error
}

func (f *fakePurchaseChecker) GetByUserAndCourse(_ context.Context, _, _ int) (*domain.Purchase, error) {
	return f.purchase, f.err
}

type fakeYooKassaGateway struct {
	called       bool
	amount       int
	currency     string
	description  string
	confirmation map[string]interface{}
	payment      *domain.Payment
	err          error
}

func (f *fakeYooKassaGateway) CreatePayment(amount int, currency string, description string, confirmation map[string]interface{}) (*domain.Payment, error) {
	f.called = true
	f.amount = amount
	f.currency = currency
	f.description = description
	f.confirmation = confirmation
	if f.err != nil {
		return nil, f.err
	}
	return f.payment, nil
}

func TestExecuteCreatesPaymentThroughGatewayAndSavesIt(t *testing.T) {
	repo := &fakePaymentRepo{}
	gateway := &fakeYooKassaGateway{
		payment: &domain.Payment{
			PaymentID:       "pay_123",
			Status:          domain.PaymentStatusPending,
			ConfirmationURL: "https://pay.example/confirm",
		},
	}
	uc := NewUseCase(
		repo,
		&fakeCourseRepo{course: domain.RestoreCourse(5, 1, "Paid course", "", false, 2500, true, "", time.Now(), "", "", "", "", "", nil)},
		&fakePurchaseChecker{err: domain.ErrPurchaseNotFound},
		gateway,
	)

	output, err := uc.Execute(context.Background(), Input{UserID: 10, CourseID: 5, ReturnURL: "https://app.example/success"})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !gateway.called {
		t.Fatal("expected YooKassa gateway mock to be called")
	}
	if gateway.amount != 2500 || gateway.currency != "RUB" {
		t.Fatalf("unexpected gateway payment data: amount=%d currency=%s", gateway.amount, gateway.currency)
	}
	if repo.created == nil {
		t.Fatal("expected payment to be saved")
	}
	if output.PaymentID != "pay_123" || output.ConfirmationURL != "https://pay.example/confirm" {
		t.Fatalf("unexpected output: %#v", output)
	}
}

func TestExecuteRejectsActiveExistingPurchase(t *testing.T) {
	gateway := &fakeYooKassaGateway{}
	uc := NewUseCase(
		&fakePaymentRepo{},
		&fakeCourseRepo{course: domain.RestoreCourse(5, 1, "Paid course", "", false, 2500, true, "", time.Now(), "", "", "", "", "", nil)},
		&fakePurchaseChecker{purchase: domain.RestorePurchase(10, 5, time.Now())},
		gateway,
	)

	_, err := uc.Execute(context.Background(), Input{UserID: 10, CourseID: 5, ReturnURL: "https://app.example/success"})
	if !errors.Is(err, domain.ErrPaymentAlreadyPaid) {
		t.Fatalf("expected ErrPaymentAlreadyPaid, got %v", err)
	}
	if gateway.called {
		t.Fatal("gateway should not be called for active existing purchase")
	}
}

func TestExecuteRejectsFreeCoursePayment(t *testing.T) {
	gateway := &fakeYooKassaGateway{}
	uc := NewUseCase(
		&fakePaymentRepo{},
		&fakeCourseRepo{course: domain.RestoreCourse(5, 1, "Free course", "", true, 0, true, "", time.Now(), "", "", "", "", "", nil)},
		&fakePurchaseChecker{err: domain.ErrPurchaseNotFound},
		gateway,
	)

	_, err := uc.Execute(context.Background(), Input{UserID: 10, CourseID: 5, ReturnURL: "https://app.example/success"})
	if err == nil {
		t.Fatal("expected error for free course payment")
	}
	if gateway.called {
		t.Fatal("gateway should not be called for free course")
	}
}
