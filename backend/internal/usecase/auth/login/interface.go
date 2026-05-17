package login

import (
	"context"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	Email    string
	Password string
}

type Output struct {
	User         *domain.User
	RefreshToken string
}

type RefreshTokenRepo interface {
	Save(ctx context.Context, userID int, token string, expiresAt time.Time) error
}

type Usecase struct {
	authenticator Authenticator
	refreshRepo   RefreshTokenRepo
}

type Authenticator interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error)
}
