package get

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	CourseID int
}

type Output struct {
	Modules []*domain.Module
}

type Usecase struct {
	reader ModuleReader
}

type ModuleReader interface {
	GetByCourseID(ctx context.Context, courseID int) ([]*domain.Module, error)
}

func NewUsecase(reader ModuleReader) (*Usecase, error) {
	if reader == nil {
		return nil, errors.New("module reader is required")
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetModulesByCourse started", zap.Int("courseID", input.CourseID))

	modules, err := u.reader.GetByCourseID(ctx, input.CourseID)
	if err != nil {
		zap.L().Error("Failed to get modules by course ID", zap.Int("courseID", input.CourseID), zap.Error(err))
		return nil, err
	}

	zap.L().Info("Modules retrieved successfully", zap.Int("courseID", input.CourseID), zap.Int("count", len(modules)))
	return &Output{Modules: modules}, nil
}
