package get

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	ModuleID int
}

type Output struct {
	Lessons []*domain.Lesson
}

type Usecase struct {
	reader LessonReader
}

func NewUsecase(reader LessonReader) (*Usecase, error) {
	if reader == nil {
		return nil, errors.New("lesson reader is required")
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	lessons, err := u.reader.GetByModuleID(ctx, input.ModuleID)
	if err != nil {
		return nil, err
	}
	return &Output{Lessons: lessons}, nil
}
