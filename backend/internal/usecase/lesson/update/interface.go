package update

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type LessonUpdater interface {
	Update(ctx context.Context, lesson *domain.Lesson) error
}

type LessonFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Lesson, error)
}
