package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseSaver interface {
	Save(ctx context.Context, course *domain.Course) error
}
