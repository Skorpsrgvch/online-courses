package approve

import (
	"context"
	"errors"
)

type Input struct {
	ReviewID int
}

type Usecase struct {
	reviewApprover ReviewApprover
}

func NewUsecase(approver ReviewApprover) (*Usecase, error) {
	if approver == nil {
		return nil, errors.New("reviewApprover is required")
	}
	return &Usecase{reviewApprover: approver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	if input.ReviewID <= 0 {
		return errors.New("invalid review ID")
	}
	return u.reviewApprover.ApproveReview(ctx, input.ReviewID)
}
