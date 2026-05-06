package get

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	ID int `json:"id"`
}

type ServiceReader interface {
	GetByID(ctx context.Context, id int) (*domain.Service, error)
}

type Usecase struct {
	reader ServiceReader
}

// NewUsecase создает новый usecase
func NewUsecase(reader ServiceReader) (*Usecase, error) {
	if reader == nil {
		return nil, nil
	}
	return &Usecase{reader: reader}, nil
}

// Execute выполняет логику
func (u *Usecase) Execute(ctx context.Context, input Input) (*domain.Service, error) {
	return u.reader.GetByID(ctx, input.ID)
}
