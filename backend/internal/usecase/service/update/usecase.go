package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
	if input.ID <= 0 {
		return domain.ErrInvalidCredentials
	}

	existing, err := u.repo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrServiceNotFound
	}

	// Обновляем поля
	existing.Title = input.Title
	existing.Price = input.Price
	existing.Description = input.Description
	existing.Duration = input.Duration

	return u.repo.Update(ctx, existing)
}
