package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Duration    int    `json:"duration_minutes"`
}

type Usecase struct {
	repo ServiceRepo
}

type ServiceRepo interface {
	Create(ctx context.Context, s *domain.Service) error
}

func NewUsecase(repo ServiceRepo) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	return &Usecase{repo: repo}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Debug("CreateService started", zap.String("title", input.Title))

	if input.Title == "" || input.Price < 0 {
		zap.L().Warn("Validation failed for service creation", zap.String("title", input.Title), zap.Int("price", input.Price))
		return domain.ErrInvalidCredentials
	}

	service := &domain.Service{
		Title:       input.Title,
		Price:       input.Price,
		Description: input.Description,
		Duration:    input.Duration,
	}

	if err := u.repo.Create(ctx, service); err != nil {
		zap.L().Error("Failed to create service", zap.Error(err))
		return err
	}

	zap.L().Info("Service created successfully", zap.Int("serviceID", service.ID))
	return nil
}
