package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ReviewSaver interface {
	CreateReview(ctx context.Context, review *domain.Review) error
}

type UserFinder interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type CourseFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}
