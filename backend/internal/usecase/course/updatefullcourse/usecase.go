package updatefullcourse

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	CourseID          int `json:"-"`
	Title             string
	Description       string
	IsPublic          bool
	Price             int
	IsActive          bool
	AuthorID          int
	CoverImageURL     string
	Contraindications string
	Recommendations   string
	TargetAudience    string
	CourseBasis       string
	ClassBasis        string
	Bonuses           []domain.BonusItem
	Modules           []ModuleInput
}

type LessonInput struct {
	ID           int
	Title        string
	Description  string
	VideoEmbedID string
	PrivateKey   *string
	Order        int
}

type ModuleInput struct {
	ID      int
	Title   string
	Order   int
	Lessons []LessonInput
}

type CourseRepository interface {
	UpdateFullWithModules(ctx context.Context, input Input) error
	ReorderModules(ctx context.Context, courseID int, orderMap map[int]int) error
	ReorderLessons(ctx context.Context, moduleID int, orderMap map[int]int) error
}

type Usecase struct {
	repo CourseRepository
}

// NewUsecase создает новый usecase и проверяет наличие репозитория
func NewUsecase(repo CourseRepository) (*Usecase, error) {
	if repo == nil {
		return nil, errors.New("course repository is required")
	}
	return &Usecase{
		repo: repo,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	zap.L().Info("Updating full course with modules", zap.Int("course_id", input.CourseID))

	if err := u.repo.UpdateFullWithModules(ctx, input); err != nil {
		zap.L().Error("Full course update failed", zap.Int("course_id", input.CourseID), zap.Error(err))
		return err
	}

	zap.L().Info("Full course updated successfully", zap.Int("course_id", input.CourseID))
	return nil
}

func (u *Usecase) ReorderModules(ctx context.Context, courseID int, moduleIDs []int) error {
	zap.L().Info("Reordering modules", zap.Int("course_id", courseID), zap.Ints("module_ids", moduleIDs))

	if len(moduleIDs) == 0 {
		return errors.New("empty module list")
	}

	orderMap := make(map[int]int)
	for i, id := range moduleIDs {
		orderMap[id] = i + 1
	}

	if err := u.repo.ReorderModules(ctx, courseID, orderMap); err != nil {
		zap.L().Error("Module reordering failed", zap.Error(err))
		return err
	}

	zap.L().Info("Modules reordered successfully")
	return nil
}

func (u *Usecase) ReorderLessons(ctx context.Context, moduleID int, lessonIDs []int) error {
	zap.L().Info("Reordering lessons", zap.Int("module_id", moduleID), zap.Ints("lesson_ids", lessonIDs))

	if len(lessonIDs) == 0 {
		return errors.New("empty lesson list")
	}

	orderMap := make(map[int]int)
	for i, id := range lessonIDs {
		orderMap[id] = i + 1
	}

	if err := u.repo.ReorderLessons(ctx, moduleID, orderMap); err != nil {
		zap.L().Error("Lesson reordering failed", zap.Error(err))
		return err
	}

	zap.L().Info("Lessons reordered successfully")
	return nil
}
