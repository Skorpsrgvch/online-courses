package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	CourseID int
	Title    string
	Order    int
}

type Usecase struct {
	moduleSaver ModuleSaver
}

type ModuleSaver interface {
	Save(ctx context.Context, module *domain.Module) error
}

func NewUsecase(moduleSaver ModuleSaver) (*Usecase, error) {
	if moduleSaver == nil {
		return nil, errors.New("moduleSaver is required")
	}
	return &Usecase{moduleSaver: moduleSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("CreateModule started", zap.Int("courseID", input.CourseID), zap.String("title", input.Title))

	module, err := domain.NewModule(input.Title, input.CourseID, input.Order)
	if err != nil {
		zap.L().Error("Failed to create module domain object", zap.Error(err))
		return err
	}

	if err := u.moduleSaver.Save(ctx, module); err != nil {
		zap.L().Error("Failed to save module", zap.Error(err))
		return err
	}

	zap.L().Info("Module created successfully", zap.Int("moduleID", module.ID), zap.Int("courseID", input.CourseID))
	return nil
}
