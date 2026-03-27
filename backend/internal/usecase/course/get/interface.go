package get

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseReader interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}

type PurchaseChecker interface {
	HasPurchased(ctx context.Context, userID, courseID int) (bool, error)
}
