package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

// Input для создания простого курса
type Input struct {
	Title             string
	Description       string
	IsPublic          bool
	Price             int
	AuthorID          int
	CoverImageURL     string
	Contraindications string
	Recommendations   string
	TargetAudience    string
	CourseBasis       string
	ClassBasis        string
	Bonuses           []domain.BonusItem
}

// ModuleInput — модуль для создания вместе с курсом
type ModuleInput struct {
	Title   string        `json:"title"`
	Order   int           `json:"order"`
	Lessons []LessonInput `json:"lessons"`
}

// LessonInput — урок для создания вместе с модулем
type LessonInput struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type Usecase struct {
	courseSaver CourseSaver
}

// CourseSaver — интерфейс для сохранения одного курса
type CourseSaver interface {
	Save(ctx context.Context, course *domain.Course) error
}

// CourseModuleSaver — интерфейс для сохранения курса с модулями (транзакция)
type CourseModuleSaver interface {
	SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []ModuleInput) error
}
