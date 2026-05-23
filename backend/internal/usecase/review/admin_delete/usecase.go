package admin_delete

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	ReviewID int
}

type Usecase struct {
	reviewDeleter ReviewDeleter
}

type ReviewDeleter interface {
	DeleteByAdmin(ctx context.Context, reviewID int) error
}

func NewUsecase(reviewDeleter ReviewDeleter) (*Usecase, error) {
	if reviewDeleter == nil {
		return nil, errors.New("reviewDeleter is required")
	}
	return &Usecase{reviewDeleter: reviewDeleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("AdminDeleteReview started", zap.Int("reviewID", input.ReviewID))

	if input.ReviewID <= 0 {
		err := errors.New("invalid review ID")
		zap.L().Warn("Validation failed", zap.Error(err))
		return err
	}

	if err := u.reviewDeleter.DeleteByAdmin(ctx, input.ReviewID); err != nil {
		zap.L().Error("Failed to delete review by admin", zap.Int("reviewID", input.ReviewID), zap.Error(err))
		return err
	}

	zap.L().Info("Review deleted by admin", zap.Int("reviewID", input.ReviewID))
	return nil
}
