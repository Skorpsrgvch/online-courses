package get

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
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

func NewUsecase(reader ServiceReader) (*Usecase, error) {
	if reader == nil {
		return nil, errors.New("reader is required")
	}
	return &Usecase{reader: reader}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*domain.Service, error) {
	zap.L().Debug("GetServiceByID started", zap.Int("serviceID", input.ID))

	service, err := u.reader.GetByID(ctx, input.ID)
	if err != nil {
		zap.L().Error("Failed to get service by ID", zap.Int("serviceID", input.ID), zap.Error(err))
		return nil, err
	}

	zap.L().Info("Service retrieved successfully", zap.Int("serviceID", input.ID))
	return service, nil
}
