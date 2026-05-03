package list

import (
	"context"
)

func NewUseCase(paymentRepo PaymentRepository) *UseCase {
	return &UseCase{
		paymentRepo: paymentRepo,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	payments, err := u.paymentRepo.ListByUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	result := make([]PaymentItem, 0, len(payments))
	for _, p := range payments {
		var paidAtStr *string
		if p.PaidAt != nil {
			s := p.PaidAt.Format("2006-01-02T15:04:05Z")
			paidAtStr = &s
		}

		result = append(result, PaymentItem{
			PaymentID:       p.PaymentID,
			CourseID:        p.CourseID,
			Amount:          p.Amount,
			Currency:        p.Currency,
			Status:          string(p.Status),
			CreatedAt:       p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			PaidAt:          paidAtStr,
			ConfirmationURL: p.ConfirmationURL,
		})
	}

	return &Output{
		Payments: result,
	}, nil
}
