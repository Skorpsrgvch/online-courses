package get

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ModuleReader interface {
	GetByCourseID(ctx context.Context, courseID int) ([]*domain.Module, error)
}
