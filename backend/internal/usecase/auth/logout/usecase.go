package logout

import (
	"context"
	"errors"

	"go.uber.org/zap"
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
	zap.L().Debug("Logout requested", zap.Int("user_id", userID))

	if err := u.repo.DeleteByUser(ctx, userID); err != nil {
		zap.L().Error("Logout failed", zap.Error(err), zap.Int("user_id", userID))
		return err
	}

	zap.L().Info("Logout successful", zap.Int("user_id", userID))
	return nil
}
