package reject

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	ReviewID int
	Reason   string
}

type Usecase struct {
	reviewDeleter ReviewDeleter
}

type ReviewDeleter interface {
	RejectReview(ctx context.Context, reviewID int, reason string) error
}

func NewUsecase(reviewDeleter ReviewDeleter) (*Usecase, error) {
	if reviewDeleter == nil {
		return nil, errors.New("reviewDeleter is required")
	}
	return &Usecase{reviewDeleter: reviewDeleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("RejectReview started", zap.Int("reviewID", input.ReviewID))

	if input.Reason == "" {
		err := errors.New("rejection reason is required")
		zap.L().Warn("Validation failed", zap.Error(err))
		return err
	}

	if err := u.reviewDeleter.RejectReview(ctx, input.ReviewID, input.Reason); err != nil {
		zap.L().Error("Failed to reject review", zap.Int("reviewID", input.ReviewID), zap.Error(err))
		return err
	}

	zap.L().Info("Review rejected", zap.Int("reviewID", input.ReviewID))
	return nil
}
