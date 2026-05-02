package list

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseLister interface {
	ListAll(ctx context.Context) ([]*domain.Course, error)
	GetAllWithInactive(ctx context.Context) ([]*domain.Course, error)
}
