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
	deleter LessonDeleter
}

type LessonDeleter interface {
	Delete(ctx context.Context, id int) error
}

func NewUsecase(deleter LessonDeleter) (*Usecase, error) {
	if deleter == nil {
		return nil, errors.New("lesson deleter is required")
	}
	return &Usecase{deleter: deleter}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("DeleteLesson started", zap.Int("lessonID", input.ID))

	if err := u.deleter.Delete(ctx, input.ID); err != nil {
		zap.L().Error("Failed to delete lesson", zap.Int("lessonID", input.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Lesson deleted successfully", zap.Int("lessonID", input.ID))
	return nil
}
