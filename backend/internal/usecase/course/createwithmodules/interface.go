package createwithmodules

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
)

type Input struct {
	Title             string               `json:"title"`
	Description       string               `json:"description"`
	IsPublic          bool                 `json:"is_public"`
	Price             int                  `json:"price"`
	AuthorID          int                  `json:"author_id"`
	CoverImageURL     string               `json:"cover_image_url"`
	Contraindications string               `json:"contraindications"`
	Recommendations   string               `json:"recommendations"`
	TargetAudience    string               `json:"target_audience"`
	CourseBasis       string               `json:"course_basis"`
	ClassBasis        string               `json:"class_basis"`
	Bonuses           []domain.BonusItem   `json:"bonuses"`
	Modules           []create.ModuleInput `json:"modules"`
}

type Usecase struct {
	courseSaver create.CourseModuleSaver
}

type CourseModuleSaver interface {
	SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []create.ModuleInput) error
}
