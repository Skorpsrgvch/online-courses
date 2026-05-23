package approve

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	ReviewID int
}

type Usecase struct {
	reviewApprover ReviewApprover
}

type ReviewApprover interface {
	ApproveReview(ctx context.Context, reviewID int) error
}

func NewUsecase(approver ReviewApprover) (*Usecase, error) {
	if approver == nil {
		return nil, errors.New("reviewApprover is required")
	}
	return &Usecase{reviewApprover: approver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("ApproveReview started", zap.Int("reviewID", input.ReviewID))

	if input.ReviewID <= 0 {
		err := errors.New("invalid review ID")
		zap.L().Warn("Validation failed", zap.Error(err))
		return err
	}

	if err := u.reviewApprover.ApproveReview(ctx, input.ReviewID); err != nil {
		zap.L().Error("Failed to approve review", zap.Int("reviewID", input.ReviewID), zap.Error(err))
		return err
	}

	zap.L().Info("Review approved", zap.Int("reviewID", input.ReviewID))
	return nil
}
