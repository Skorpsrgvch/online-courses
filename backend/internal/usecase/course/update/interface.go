package update

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseUpdater interface {
	Update(ctx context.Context, course *domain.Course) error
}

type CourseFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}
