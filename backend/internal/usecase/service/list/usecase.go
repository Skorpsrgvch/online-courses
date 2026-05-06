package list

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ServiceReader interface {
	GetAll(ctx context.Context) ([]*domain.Service, error)
}

type Usecase struct {
	reader ServiceReader
}

func NewUsecase(reader ServiceReader) (*Usecase, error) {
	if reader == nil {
		return nil, nil
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context) ([]*domain.Service, error) {
	return u.reader.GetAll(ctx)
}
