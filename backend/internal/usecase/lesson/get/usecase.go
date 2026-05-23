package get

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	ModuleID int
}

type Output struct {
	Lessons []*domain.Lesson
}

type Usecase struct {
	reader LessonReader
}

type LessonReader interface {
	GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error)
}

func NewUsecase(reader LessonReader) (*Usecase, error) {
	if reader == nil {
		return nil, errors.New("lesson reader is required")
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetLessonsByModule started", zap.Int("moduleID", input.ModuleID))

	lessons, err := u.reader.GetByModuleID(ctx, input.ModuleID)
	if err != nil {
		zap.L().Error("Failed to get lessons by module ID", zap.Int("moduleID", input.ModuleID), zap.Error(err))
		return nil, err
	}

	zap.L().Info("Lessons retrieved successfully", zap.Int("moduleID", input.ModuleID), zap.Int("count", len(lessons)))
	return &Output{Lessons: lessons}, nil
}
