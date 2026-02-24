package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	CourseID int
	Title    string
	Order    int
}

type Usecase struct {
	moduleSaver ModuleSaver
}

func NewUsecase(moduleSaver ModuleSaver) (*Usecase, error) {
	if moduleSaver == nil {
		return nil, errors.New("moduleSaver is required")
	}
	return &Usecase{moduleSaver: moduleSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	module, err := domain.NewModule(input.Title, input.CourseID, input.Order)
	if err != nil {
		return err
	}
	return u.moduleSaver.Save(ctx, module)
}
