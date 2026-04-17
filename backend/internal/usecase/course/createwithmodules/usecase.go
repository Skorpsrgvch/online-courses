package createwithmodules

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
)

type Input struct {
	Title         string
	Description   string
	IsPublic      bool
	Price         int
	AuthorID      int
	CoverImageURL string
	Modules       []create.ModuleInput
}

type Usecase struct {
	courseSaver CourseModuleSaver
}

func NewUsecase(courseSaver CourseModuleSaver) (*Usecase, error) {
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
	return u.courseSaver.SaveCourseWithModules(ctx, course, input.Modules)
}
