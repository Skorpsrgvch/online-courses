package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
}

type ModuleInput struct {
	Title   string        `json:"title"`
	Order   int           `json:"order"`
	Lessons []LessonInput `json:"lessons"`
}

type LessonInput struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type CourseSaver interface {
	Save(ctx context.Context, course *domain.Course) error
}

type CourseModuleSaver interface {
	SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []ModuleInput) error
}

type Usecase struct {
	courseSaver CourseSaver
}

func NewUsecase(courseSaver CourseSaver) (*Usecase, error) {
	if courseSaver == nil {
		return nil, errors.New("courseSaver is required")
	}
	return &Usecase{
		courseSaver: courseSaver,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Info("Creating new course",
		zap.String("title", input.Title),
		zap.Int("author_id", input.AuthorID),
		zap.Int("price", input.Price),
	)

	bonuses := input.Bonuses
	if bonuses == nil {
		bonuses = []domain.BonusItem{}
	}

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

	course.Bonuses = bonuses

	if err := u.courseSaver.Save(ctx, course); err != nil {
		zap.L().Error("Failed to save course", zap.Error(err))
		return err
	}

	zap.L().Info("Course created successfully", zap.Int("course_id", course.ID))
	return nil
}
