package refresh

import (
	"context"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	RefreshToken string `json:"refresh_token"`
	UserID       int
}

type Output struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
