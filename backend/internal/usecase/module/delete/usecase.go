package delete

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	ID int
}

type Usecase struct {
	deleter ModuleDeleter
}

type ModuleDeleter interface {
	Delete(ctx context.Context, id int) error
}

func NewUsecase(deleter ModuleDeleter) (*Usecase, error) {
	if deleter == nil {
		return nil, errors.New("module deleter is required")
	}
	return &Usecase{deleter: deleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("DeleteModule started", zap.Int("moduleID", input.ID))

	if err := u.deleter.Delete(ctx, input.ID); err != nil {
		zap.L().Error("Failed to delete module", zap.Int("moduleID", input.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Module deleted successfully", zap.Int("moduleID", input.ID))
	return nil
}
