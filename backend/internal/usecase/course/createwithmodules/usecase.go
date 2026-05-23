package createwithmodules

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	"go.uber.org/zap"
)

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
	Modules           []create.ModuleInput
}

type CourseModuleSaver interface {
	SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []create.ModuleInput) error
}

type Usecase struct {
	courseSaver CourseModuleSaver
}

func NewUsecase(courseSaver CourseModuleSaver) (*Usecase, error) {
	if courseSaver == nil {
		return nil, errors.New("courseSaver is required")
	}
	return &Usecase{
		courseSaver: courseSaver,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Info("Creating course with modules",
		zap.String("title", input.Title),
		zap.Int("modules_count", len(input.Modules)),
	)

	course, err := domain.NewCourse(
		input.Title,
		input.Description,
		input.IsPublic,
		input.Price,
		input.AuthorID,
		input.CoverImageURL,
		input.Contraindications,
		input.Recommendations,
		input.TargetAudience,
		input.CourseBasis,
		input.ClassBasis,
	)
	if err != nil {
		zap.L().Warn("Course validation failed", zap.Error(err))
		return err
	}

	if input.Bonuses == nil {
		course.Bonuses = []domain.BonusItem{}
	} else {
		course.Bonuses = input.Bonuses
	}

	if err := u.courseSaver.SaveCourseWithModules(ctx, course, input.Modules); err != nil {
		zap.L().Error("Failed to save course with modules", zap.Error(err))
		return err
	}

	zap.L().Info("Course with modules created successfully", zap.Int("course_id", course.ID))
	return nil
}
