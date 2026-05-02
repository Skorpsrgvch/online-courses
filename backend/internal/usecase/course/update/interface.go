package update

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
	Update(ctx context.Context, course *domain.Course) error
	UpdateStatus(ctx context.Context, id int, isActive bool) error
}

type Usecase struct {
	repo CourseRepository
}
