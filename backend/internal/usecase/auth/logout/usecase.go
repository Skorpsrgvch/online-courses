package logout

import (
	"context"
	"errors"
)

type RefreshTokenRepo interface {
	DeleteByUser(ctx context.Context, userID int) error
}

type Usecase struct {
	repo RefreshTokenRepo
}

func NewUsecase(repo RefreshTokenRepo) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("repo required")
	}
	return &Usecase{repo: repo}, nil
}

func (u *Usecase) Execute(ctx context.Context, userID int) error {
	return u.repo.DeleteByUser(ctx, userID)
}
