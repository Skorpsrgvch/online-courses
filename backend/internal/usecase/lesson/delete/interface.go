package delete

import "context"

type LessonDeleter interface {
	Delete(ctx context.Context, lessonID int) error
}
