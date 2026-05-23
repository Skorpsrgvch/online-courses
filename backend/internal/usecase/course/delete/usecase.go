package delete

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Input struct {
	CourseID int
}

type CourseDeleter interface {
	SetInactive(ctx context.Context, courseID int) error
}

type Usecase struct {
	deleter CourseDeleter
}

func NewUsecase(deleter CourseDeleter) (*Usecase, error) {
	if deleter == nil {
		return nil, errors.New("deleter is required")
	}
	return &Usecase{
		deleter: deleter,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Info("Deleting course (setting inactive)", zap.Int("course_id", input.CourseID))

	if err := u.deleter.SetInactive(ctx, input.CourseID); err != nil {
		zap.L().Error("Failed to delete course", zap.Int("course_id", input.CourseID), zap.Error(err))
		return err
	}

	zap.L().Info("Course deleted successfully", zap.Int("course_id", input.CourseID))
	return nil
}
