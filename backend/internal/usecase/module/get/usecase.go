package get

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	CourseID int
}

type Output struct {
	Modules []*domain.Module
}

type Usecase struct {
	reader ModuleReader
}

func NewUsecase(reader ModuleReader) (*Usecase, error) {
	if reader == nil {
		return nil, errors.New("module reader is required")
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	modules, err := u.reader.GetByCourseID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}
	return &Output{Modules: modules}, nil
}
