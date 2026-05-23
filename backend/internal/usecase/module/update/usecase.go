package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	ID    int
	Title string
	Order int
}

type Usecase struct {
	updater ModuleUpdater
	finder  ModuleFinder
}

type ModuleUpdater interface {
	Update(ctx context.Context, module *domain.Module) error
}

type ModuleFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Module, error)
}

func NewUsecase(updater ModuleUpdater, finder ModuleFinder) (*Usecase, error) {
	if updater == nil || finder == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{updater: updater, finder: finder}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("UpdateModule started", zap.Int("moduleID", input.ID))

	existing, err := u.finder.GetByID(ctx, input.ID)
	if err != nil {
		zap.L().Error("Failed to find existing module", zap.Int("moduleID", input.ID), zap.Error(err))
		return err
	}

	updated := domain.RestoreModule(input.ID, existing.CourseID, input.Order, input.Title)

	if err := u.updater.Update(ctx, updated); err != nil {
		zap.L().Error("Failed to update module", zap.Int("moduleID", input.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Module updated successfully", zap.Int("moduleID", input.ID))
	return nil
}
