package mark

import "context"

type ProgressSaver interface {
	MarkCompleted(ctx context.Context, userID, lessonID int) error
}
