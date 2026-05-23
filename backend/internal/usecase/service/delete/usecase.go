package delete

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	ID int
}

type Usecase struct {
	repo ServiceRepo
}

type ServiceRepo interface {
	Delete(ctx context.Context, id int) error
}

func NewUsecase(repo ServiceRepo) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	return &Usecase{repo: repo}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("DeleteService started", zap.Int("serviceID", input.ID))

	if input.ID <= 0 {
		zap.L().Warn("Invalid service ID for deletion", zap.Int("id", input.ID))
		return domain.ErrInvalidCredentials
	}

	if err := u.repo.Delete(ctx, input.ID); err != nil {
		zap.L().Error("Failed to delete service", zap.Int("serviceID", input.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Service deleted successfully", zap.Int("serviceID", input.ID))
	return nil
}
