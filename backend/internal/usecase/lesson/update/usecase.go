package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	ID           int
	Title        string
	Description  string
	VideoEmbedID string
	PrivateKey   *string
	Order        int
}

type Usecase struct {
	updater LessonUpdater
	finder  LessonFinder
}

type LessonUpdater interface {
	Update(ctx context.Context, lesson *domain.Lesson) error
}

type LessonFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Lesson, error)
}

func NewUsecase(updater LessonUpdater, finder LessonFinder) (*Usecase, error) {
	if updater == nil || finder == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{updater: updater, finder: finder}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("UpdateLesson started", zap.Int("lessonID", input.ID))

	existing, err := u.finder.GetByID(ctx, input.ID)
	if err != nil {
		zap.L().Error("Failed to find existing lesson", zap.Int("lessonID", input.ID), zap.Error(err))
		return err
	}

	updated := domain.RestoreLesson(
		input.ID,
		existing.ModuleID,
		input.Order,
		input.Title,
		input.Description,
		input.VideoEmbedID,
		input.PrivateKey,
	)

	if err := u.updater.Update(ctx, updated); err != nil {
		zap.L().Error("Failed to update lesson", zap.Int("lessonID", input.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Lesson updated successfully", zap.Int("lessonID", input.ID))
	return nil
}
