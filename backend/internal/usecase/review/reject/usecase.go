package reject

import (
	"context"
	"errors"
)

type Input struct {
	ReviewID int
}

type ReviewDeleter interface {
	RejectReview(ctx context.Context, reviewID int) error
}

type Usecase struct {
	reviewDeleter ReviewDeleter
}

func NewUsecase(reviewDeleter ReviewDeleter) (*Usecase, error) {
	if reviewDeleter == nil {
		return nil, errors.New("reviewDeleter is required")
	}
	return &Usecase{reviewDeleter: reviewDeleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	return u.reviewDeleter.RejectReview(ctx, input.ReviewID)
}
