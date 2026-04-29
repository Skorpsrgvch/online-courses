package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

func NewUsecase(lessonSaver LessonSaver) (*Usecase, error) {
	if lessonSaver == nil {
		return nil, errors.New("lessonSaver is required")
	}
	return &Usecase{lessonSaver: lessonSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	lesson, err := domain.NewLesson(
		input.Title,
		input.Description,
		input.VideoEmbedID,
		input.ModuleID,
		input.Order,
		input.PrivateKey,
	)
	if err != nil {
		return err
	}
	return u.lessonSaver.Save(ctx, lesson)
}
