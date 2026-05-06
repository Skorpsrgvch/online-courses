package delete

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
	if input.ID <= 0 {
		return domain.ErrInvalidCredentials
	}
	return u.repo.Delete(ctx, input.ID)
}
