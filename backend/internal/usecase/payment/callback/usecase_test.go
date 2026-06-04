package callback

import (
	"context"
	"testing"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type fakeCallbackPaymentRepo struct {
	payment *domain.Payment
	updated *domain.Payment
	err     error
}

func (f *fakeCallbackPaymentRepo) GetByID(_ context.Context, _ string) (*domain.Payment, error) {
	return f.payment, f.err
}

func (f *fakeCallbackPaymentRepo) Update(_ context.Context, payment *domain.Payment) error {
	f.updated = payment
	return nil
}

type fakeCallbackPurchaseRepo struct {
	existing         *domain.Purchase
	createCalled    bool
	updateDateCalled bool
}

func (f *fakeCallbackPurchaseRepo) GetByUserAndCourse(context.Context, int, int) (*domain.Purchase, error) {
	if f.existing == nil {
		return nil, domain.ErrPurchaseNotFound
	}
	return f.existing, nil
}

func (f *fakeCallbackPurchaseRepo) Create(context.Context, int, int, string) error {
	f.createCalled = true
	return nil
}

func (f *fakeCallbackPurchaseRepo) UpdatePurchaseDate(context.Context, int, int, time.Time) error {
	f.updateDateCalled = true
	return nil
}

type fakeYooKassaVerifier struct {
	status string
	paid   bool
	err    error
}

func (f *fakeYooKassaVerifier) GetPaymentStatus(string) (string, bool, error) {
	return f.status, f.paid, f.err
}

func TestExecuteCreatesPurchaseForSucceededPayment(t *testing.T) {
	payment := domain.NewPayment("pay_1", 10, 5, 2500, "RUB", "Course payment")
	purchaseRepo := &fakeCallbackPurchaseRepo{}
	uc := NewUseCase(
		&fakeCallbackPaymentRepo{payment: payment},
		purchaseRepo,
		&fakeYooKassaVerifier{status: "succeeded", paid: true},
	)

	output, err := uc.Execute(context.Background(), Input{PaymentID: "pay_1", Status: "succeeded"})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if output == nil || !output.Processed || output.Status != "succeeded" {
		t.Fatalf("unexpected output: %#v", output)
	}
	if !purchaseRepo.createCalled {
		t.Fatal("expected purchase to be created")
	}
	if purchaseRepo.updateDateCalled {
		t.Fatal("purchase date should not be updated when purchase does not exist")
	}
}

func TestExecuteRenewsExistingPurchaseForSucceededPayment(t *testing.T) {
	payment := domain.NewPayment("pay_2", 10, 5, 2500, "RUB", "Course payment")
	purchaseRepo := &fakeCallbackPurchaseRepo{
		existing: domain.RestorePurchase(10, 5, time.Now().Add(-400*24*time.Hour)),
	}
	uc := NewUseCase(
		&fakeCallbackPaymentRepo{payment: payment},
		purchaseRepo,
		&fakeYooKassaVerifier{status: "waiting_for_capture", paid: true},
	)

	output, err := uc.Execute(context.Background(), Input{PaymentID: "pay_2", Status: "waiting_for_capture"})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if output == nil || !output.Processed || output.Status != "succeeded" {
		t.Fatalf("unexpected output: %#v", output)
	}
	if !purchaseRepo.updateDateCalled {
		t.Fatal("expected existing purchase date to be renewed")
	}
	if purchaseRepo.createCalled {
		t.Fatal("new purchase should not be created when purchase exists")
	}
}
