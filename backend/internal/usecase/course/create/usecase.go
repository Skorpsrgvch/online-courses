package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	Title         string
	Description   string
	IsPublic      bool
	Price         int
	AuthorID      int
	CoverImageURL string
}

type Usecase struct {
	courseSaver CourseSaver
}

func NewUsecase(courseSaver CourseSaver) (*Usecase, error) {
	if courseSaver == nil {
		return nil, errors.New("courseSaver is required")
	}
	return &Usecase{courseSaver: courseSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	course, err := domain.NewCourse(
		input.Title,
		input.Description,
		input.IsPublic,
		input.Price,
		input.AuthorID,
		input.CoverImageURL,
	)
	if err != nil {
		return err
	}
	return u.courseSaver.Save(ctx, course)
}
