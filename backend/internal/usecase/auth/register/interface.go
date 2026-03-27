package register

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *domain.User, passwordHash string) error
}
