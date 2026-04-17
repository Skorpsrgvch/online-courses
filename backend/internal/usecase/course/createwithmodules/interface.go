package createwithmodules

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
)

type CourseModuleSaver interface {
	SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []create.ModuleInput) error
}
