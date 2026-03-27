package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	ID           int
	Title        string
	Description  string
	VideoEmbedID string
	Order        int
}

type Usecase struct {
	updater LessonUpdater
	finder  LessonFinder
}

func NewUsecase(updater LessonUpdater, finder LessonFinder) (*Usecase, error) {
	if updater == nil || finder == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{updater: updater, finder: finder}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	existing, err := u.finder.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	updated := domain.RestoreLesson(
		input.ID, existing.ModuleID, input.Order,
		input.Title, input.Description, input.VideoEmbedID,
	)
	return u.updater.Update(ctx, updated)
}
