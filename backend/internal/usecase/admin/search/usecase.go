package search

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
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

func NewUsecase(repo UserRepository) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("repository is required")
	}
	return &Usecase{repo: repo}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) ([]OutputUser, error) {
	zap.L().Debug("Search users started", zap.String("query", input.EmailQuery), zap.Int("limit", input.Limit))

	if input.EmailQuery == "" {
		return nil, errors.New("поисковый запрос пуст")
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	users, err := u.repo.SearchByEmail(ctx, input.EmailQuery, limit)
	if err != nil {
		zap.L().Error("Search failed", zap.Error(err))
		return nil, err
	}

	if len(users) == 0 {
		zap.L().Debug("No users found")
		return []OutputUser{}, nil
	}

	result := make([]OutputUser, 0, len(users))
	for _, user := range users {
		result = append(result, OutputUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	zap.L().Info("Search completed", zap.Int("found_count", len(result)))
	return result, nil
}
