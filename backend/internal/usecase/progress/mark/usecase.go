package mark

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	UserID   int
	LessonID int
}

type Usecase struct {
	progressSaver ProgressSaver
}

type ProgressSaver interface {
	MarkCompleted(ctx context.Context, userID, lessonID int) error
}

func NewUsecase(progressSaver ProgressSaver) (*Usecase, error) {
	if progressSaver == nil {
		return nil, errors.New("progressSaver is required")
	}
	return &Usecase{progressSaver: progressSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("MarkLessonCompleted started", zap.Int("userID", input.UserID), zap.Int("lessonID", input.LessonID))

	if input.UserID <= 0 || input.LessonID <= 0 {
		err := errors.New("invalid user or lesson ID")
		zap.L().Warn("Validation failed", zap.Error(err))
		return domain.ErrInvalidInput
	}

	if err := u.progressSaver.MarkCompleted(ctx, input.UserID, input.LessonID); err != nil {
		zap.L().Error("Failed to mark lesson as completed", zap.Int("userID", input.UserID), zap.Int("lessonID", input.LessonID), zap.Error(err))
		return err
	}

	zap.L().Info("Lesson marked as completed", zap.Int("userID", input.UserID), zap.Int("lessonID", input.LessonID))
	return nil
}
