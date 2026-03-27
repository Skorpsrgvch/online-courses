package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	ID          int
	Title       string
	Description string
	IsPublic    bool
	Price       int
	IsActive    bool
	AuthorID    int
}

type Usecase struct {
	courseUpdater CourseUpdater
	courseFinder  CourseFinder
}

func NewUsecase(updater CourseUpdater, finder CourseFinder) (*Usecase, error) {
	if updater == nil || finder == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{
		courseUpdater: updater,
		courseFinder:  finder,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	// Проверяем существование
	_, err := u.courseFinder.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	// Обновляем поля
	updated := domain.RestoreCourse(
		input.ID,
		input.Title,
		input.Description,
		input.IsPublic,
		input.Price,
		input.AuthorID,
		input.IsActive,
	)

	return u.courseUpdater.Update(ctx, updated)
}
