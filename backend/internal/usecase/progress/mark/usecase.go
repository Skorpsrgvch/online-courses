package mark

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	UserID   int
	LessonID int
}

type Usecase struct {
	progressSaver ProgressSaver
}

func NewUsecase(progressSaver ProgressSaver) (*Usecase, error) {
	if progressSaver == nil {
		return nil, errors.New("progressSaver is required")
	}
	return &Usecase{progressSaver: progressSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	if input.UserID <= 0 || input.LessonID <= 0 {
		return domain.ErrInvalidCredentials // или новая ошибка
	}
	return u.progressSaver.MarkCompleted(ctx, input.UserID, input.LessonID)
}
