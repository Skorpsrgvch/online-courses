package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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

func NewUsecase(updater ModuleUpdater, finder ModuleFinder) (*Usecase, error) {
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
	updated := domain.RestoreModule(input.ID, existing.CourseID, input.Order, input.Title)
	return u.updater.Update(ctx, updated)
}
