package login

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Authenticator interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error)
}
