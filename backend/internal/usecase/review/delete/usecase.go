package delete

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	ReviewID int
}

type Output struct{}

type Usecase struct {
	reviewDeleter ReviewDeleter
}

type ReviewDeleter interface {
	DeleteReview(ctx context.Context, reviewID int) error
}

func NewUsecase(reviewDeleter ReviewDeleter) (*Usecase, error) {
	if reviewDeleter == nil {
		return nil, errors.New("reviewDeleter is required")
	}
	return &Usecase{reviewDeleter: reviewDeleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	zap.L().Debug("DeleteReview started", zap.Int("reviewID", input.ReviewID))

	if input.ReviewID <= 0 {
		err := errors.New("invalid review ID")
		zap.L().Warn("Validation failed", zap.Error(err))
		return Output{}, err
	}

	if err := u.reviewDeleter.DeleteReview(ctx, input.ReviewID); err != nil {
		zap.L().Error("Failed to delete review", zap.Int("reviewID", input.ReviewID), zap.Error(err))
		return Output{}, err
	}

	zap.L().Info("Review deleted", zap.Int("reviewID", input.ReviewID))
	return Output{}, nil
}
