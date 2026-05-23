package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	ModuleID     int
	Title        string
	Description  string
	VideoEmbedID string
	PrivateKey   *string
	Order        int
}

type Usecase struct {
	lessonSaver LessonSaver
}

type LessonSaver interface {
	Save(ctx context.Context, lesson *domain.Lesson) error
}

func NewUsecase(lessonSaver LessonSaver) (*Usecase, error) {
	if lessonSaver == nil {
		return nil, errors.New("lessonSaver is required")
	}
	return &Usecase{lessonSaver: lessonSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("CreateLesson started", zap.Int("moduleID", input.ModuleID), zap.String("title", input.Title))

	lesson, err := domain.NewLesson(
		input.Title,
		input.Description,
		input.VideoEmbedID,
		input.ModuleID,
		input.Order,
		input.PrivateKey,
	)
	if err != nil {
		zap.L().Error("Failed to create lesson domain object", zap.Error(err))
		return err
	}

	if err := u.lessonSaver.Save(ctx, lesson); err != nil {
		zap.L().Error("Failed to save lesson", zap.Error(err))
		return err
	}

	zap.L().Info("Lesson created successfully", zap.Int("lessonID", lesson.ID), zap.Int("moduleID", input.ModuleID))
	return nil
}
