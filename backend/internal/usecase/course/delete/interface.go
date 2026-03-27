package delete

import "context"

type CourseDeleter interface {
	SetInactive(ctx context.Context, courseID int) error
}
