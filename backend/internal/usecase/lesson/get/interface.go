package get

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type LessonReader interface {
	GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error)
}
