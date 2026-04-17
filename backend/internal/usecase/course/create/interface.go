package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

// CourseSaver — интерфейс для сохранения одного курса
type CourseSaver interface {
	Save(ctx context.Context, course *domain.Course) error
}

// ModuleInput — модуль для создания вместе с курсом
type ModuleInput struct {
	Title   string        `json:"title"`
	Order   int           `json:"order"`
	Lessons []LessonInput `json:"lessons"`
}

// LessonInput — урок для создания вместе с модулем
type LessonInput struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	LessonType     domain.LessonType `json:"lesson_type"`
	VideoEmbedID   string            `json:"video_embed_id"`
	ArticleContent string            `json:"article_content"`
	Order          int               `json:"order"`
}

// CourseModuleSaver — интерфейс для сохранения курса с модулями (транзакция)
type CourseModuleSaver interface {
	SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []ModuleInput) error
}
