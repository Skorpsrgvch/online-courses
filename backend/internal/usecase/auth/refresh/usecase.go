package refresh

import (
	"context"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/middleware"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Output struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenRepo interface {
	Validate(ctx context.Context, userID int, token string) error
	Save(ctx context.Context, userID int, token string, expiresAt time.Time) error
	DeleteByUser(ctx context.Context, userID int) error
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type Usecase struct {
	repo     RefreshTokenRepo
	userRepo UserRepo
}

func NewUsecase(repo RefreshTokenRepo, userRepo UserRepo) (*Usecase, error) {
	if repo == nil || userRepo == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{repo: repo, userRepo: userRepo}, nil
}

func (u *Usecase) Execute(ctx context.Context, refreshToken string) (*Output, error) {
	zap.L().Debug("Token refresh started")

	claims, err := middleware.ParseToken(refreshToken)
	if err != nil {
		zap.L().Warn("Invalid refresh token signature", zap.Error(err))
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().UTC().After(claims.ExpiresAt.Time) {
		zap.L().Warn("Refresh token expired", zap.Int("user_id", claims.UserID))
		return nil, errors.New("refresh token expired")
	}

	if err := u.repo.Validate(ctx, claims.UserID, refreshToken); err != nil {
		zap.L().Warn("Token not found in DB or revoked", zap.Int("user_id", claims.UserID))
		return nil, errors.New("token revoked or not found")
	}

	user, err := u.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			zap.L().Error("User not found during refresh", zap.Int("user_id", claims.UserID))
			return nil, errors.New("user not found")
		}
		zap.L().Error("DB error during refresh", zap.Error(err))
		return nil, err
	}

	newAccess, newRefresh, err := middleware.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		zap.L().Error("Token generation failed", zap.Error(err))
		return nil, err
	}

	// Rotational strategy: удаляем старый, сохраняем новый
	if err := u.repo.DeleteByUser(ctx, user.ID); err != nil {
		zap.L().Error("Failed to delete old token", zap.Error(err))
	}

	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
	if err := u.repo.Save(ctx, user.ID, newRefresh, expiresAt); err != nil {
		zap.L().Error("Failed to save new token", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Tokens refreshed successfully", zap.Int("user_id", user.ID))
	return &Output{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
