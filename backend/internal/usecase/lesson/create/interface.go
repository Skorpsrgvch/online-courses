package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type LessonSaver interface {
	Save(ctx context.Context, lesson *domain.Lesson) error
}
