package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
	if input.Title == "" || input.Price < 0 {
		return domain.ErrInvalidCredentials
	}

	service := &domain.Service{
		Title:       input.Title,
		Price:       input.Price,
		Description: input.Description,
		Duration:    input.Duration,
	}

	return u.repo.Create(ctx, service)
}
