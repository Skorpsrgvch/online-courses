package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	ModuleID     int
	Title        string
	Description  string
	VideoEmbedID string
	PrivateKey   *string
	Order        int
}

type Usecase struct {
	lessonSaver LessonSaver
}

type LessonSaver interface {
	Save(ctx context.Context, lesson *domain.Lesson) error
}
