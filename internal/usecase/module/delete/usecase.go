package delete

import (
	"context"
	"errors"
)

type Input struct {
	ID int
}

type Usecase struct {
	deleter ModuleDeleter
}

func NewUsecase(deleter ModuleDeleter) (*Usecase, error) {
	if deleter == nil {
		return nil, errors.New("module deleter is required")
	}
	return &Usecase{deleter: deleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	return u.deleter.Delete(ctx, input.ID)
}
