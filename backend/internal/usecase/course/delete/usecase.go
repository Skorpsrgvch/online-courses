package delete

import (
	"context"
	"errors"
)

type Input struct {
	CourseID int
}

type Usecase struct {
	deleter CourseDeleter
}

func NewUsecase(deleter CourseDeleter) (*Usecase, error) {
	if deleter == nil {
		return nil, errors.New("deleter is required")
	}
	return &Usecase{deleter: deleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	return u.deleter.SetInactive(ctx, input.CourseID)
}
