package delete

import (
	"context"
	"errors"
)

type Input struct {
	ID int
}

type Usecase struct {
	deleter LessonDeleter
}

func NewUsecase(deleter LessonDeleter) (*Usecase, error) {
	if deleter == nil {
		return nil, errors.New("lesson deleter is required")
	}
	return &Usecase{deleter: deleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	return u.deleter.Delete(ctx, input.ID)
}
