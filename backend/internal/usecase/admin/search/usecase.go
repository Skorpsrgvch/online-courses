package search

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	EmailQuery string `json:"email"`
	Limit      int    `json:"limit"`
}

type OutputUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserRepository interface {
	SearchByEmail(ctx context.Context, query string, limit int) ([]*domain.User, error)
}

type Usecase struct {
	repo UserRepository
}

func NewUsecase(repo UserRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(ctx context.Context, input Input) ([]OutputUser, error) {
	if input.Limit <= 0 {
		input.Limit = 10
	}

	users, err := u.repo.SearchByEmail(ctx, input.EmailQuery, input.Limit)
	if err != nil {
		return nil, err
	}

	var result []OutputUser
	for _, user := range users {
		result = append(result, OutputUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}
	return result, nil
}
