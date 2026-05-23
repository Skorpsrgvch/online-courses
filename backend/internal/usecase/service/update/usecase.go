package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	ID          int    `json:"-"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Duration    int    `json:"duration_minutes"`
}

type Usecase struct {
	repo ServiceRepo
}

type ServiceRepo interface {
	GetByID(ctx context.Context, id int) (*domain.Service, error)
	Update(ctx context.Context, s *domain.Service) error
}

func NewUsecase(repo ServiceRepo) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	return &Usecase{repo: repo}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("UpdateService started", zap.Int("serviceID", input.ID))

	if input.ID <= 0 {
		zap.L().Warn("Invalid service ID for update", zap.Int("id", input.ID))
		return domain.ErrInvalidCredentials
	}

	existing, err := u.repo.GetByID(ctx, input.ID)
	if err != nil {
		zap.L().Error("Failed to find existing service", zap.Int("serviceID", input.ID), zap.Error(err))
		return err
	}
	if existing == nil {
		zap.L().Warn("Service not found for update", zap.Int("serviceID", input.ID))
		return domain.ErrServiceNotFound
	}

	existing.Title = input.Title
	existing.Price = input.Price
	existing.Description = input.Description
	existing.Duration = input.Duration

	if err := u.repo.Update(ctx, existing); err != nil {
		zap.L().Error("Failed to update service", zap.Int("serviceID", input.ID), zap.Error(err))
		return err
	}

	zap.L().Info("Service updated successfully", zap.Int("serviceID", input.ID))
	return nil
}
