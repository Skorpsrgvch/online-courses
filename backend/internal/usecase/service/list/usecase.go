package list

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type ServiceReader interface {
	GetAll(ctx context.Context) ([]*domain.Service, error)
}

type Usecase struct {
	reader ServiceReader
}

func NewUsecase(reader ServiceReader) (*Usecase, error) {
	if reader == nil {
		return nil, errors.New("reader is required")
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context) ([]*domain.Service, error) {
	zap.L().Debug("ListServices started")

	services, err := u.reader.GetAll(ctx)
	if err != nil {
		zap.L().Error("Failed to get all services", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Services list retrieved successfully", zap.Int("count", len(services)))
	return services, nil
}
