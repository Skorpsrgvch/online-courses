package getfull

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	CourseID int
	UserID   int
	Role     string
}

type LessonOutput struct {
	ID           int     `json:"id"`
	ModuleID     int     `json:"module_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
	IsCompleted  bool    `json:"is_completed"`
}

type ModuleOutput struct {
	ID       int            `json:"id"`
	CourseID int            `json:"course_id"`
	Title    string         `json:"title"`
	Order    int            `json:"order"`
	Lessons  []LessonOutput `json:"lessons"`
}

type Output struct {
	Course          *domain.Course `json:"course"`
	Modules         []ModuleOutput `json:"modules"`
	IsPurchased     bool           `json:"is_purchased"`
	ProgressPercent int            `json:"progress_percent"`
}

type CourseReader interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}

type ModuleReader interface {
	GetByCourseID(ctx context.Context, courseID int) ([]*domain.Module, error)
}

type LessonReader interface {
	GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error)
}

type PurchaseChecker interface {
	HasPurchased(ctx context.Context, userID, courseID int) (bool, error)
}

type ProgressReader interface {
	GetCourseProgress(ctx context.Context, userID, courseID int) (completed int, total int, err error)
	IsLessonCompleted(ctx context.Context, userID, lessonID int) (bool, error)
}

type Usecase struct {
	courseReader    CourseReader
	moduleReader    ModuleReader
	lessonReader    LessonReader
	purchaseChecker PurchaseChecker
	progressReader  ProgressReader
}

func NewUsecase(courseReader CourseReader, moduleReader ModuleReader, lessonReader LessonReader, purchaseChecker PurchaseChecker, pr ProgressReader) (*Usecase, error) {
	if courseReader == nil || moduleReader == nil || lessonReader == nil || purchaseChecker == nil {
		return nil, errors.New("all dependencies required")
	}
	return &Usecase{
		courseReader:    courseReader,
		moduleReader:    moduleReader,
		lessonReader:    lessonReader,
		purchaseChecker: purchaseChecker,
		progressReader:  pr,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetFullCourse started", zap.Int("courseID", input.CourseID), zap.Int("userID", input.UserID))

	course, err := u.courseReader.GetByID(ctx, input.CourseID)
	if err != nil {
		zap.L().Error("Failed to get course by ID", zap.Int("courseID", input.CourseID), zap.Error(err))
		return nil, err
	}

	hasFullAccess := false
	isPurchased := false

	// 1. Проверяем покупку в БД
	if input.UserID > 0 {
		purchased, err := u.purchaseChecker.HasPurchased(ctx, input.UserID, course.ID)
		if err != nil {
			zap.L().Error("Failed to check purchase status", zap.Int("userID", input.UserID), zap.Int("courseID", course.ID), zap.Error(err))
			return nil, err
		}

		if purchased {
			isPurchased = true
			hasFullAccess = true
			zap.L().Debug("User has purchased the course", zap.Int("userID", input.UserID))
		}
	}

	// 2. Если покупки нет, проверяем права на доступ (бесплатный курс или админ)
	if !hasFullAccess {
		if course.IsPublic {
			hasFullAccess = true
			zap.L().Debug("Access granted via public course flag", zap.Int("courseID", course.ID))
		} else if input.Role == "admin" {
			hasFullAccess = true
			zap.L().Debug("Access granted via admin role", zap.Int("userID", input.UserID))
		} else {
			zap.L().Debug("No full access granted", zap.Bool("isPublic", course.IsPublic), zap.String("role", input.Role))
		}
	}

	dbModules, err := u.moduleReader.GetByCourseID(ctx, course.ID)
	if err != nil {
		zap.L().Error("Failed to get modules", zap.Int("courseID", course.ID), zap.Error(err))
		return nil, err
	}

	modulesOut := make([]ModuleOutput, 0, len(dbModules))

	for _, m := range dbModules {
		dbLessons, err := u.lessonReader.GetByModuleID(ctx, m.ID)
		if err != nil {
			zap.L().Error("Failed to get lessons for module", zap.Int("moduleID", m.ID), zap.Error(err))
			return nil, err
		}

		lessonsOut := make([]LessonOutput, 0, len(dbLessons))

		for _, l := range dbLessons {
			lessonOut := LessonOutput{
				ID:          l.ID,
				ModuleID:    l.ModuleID,
				Title:       l.Title,
				Description: l.Description,
				Order:       l.Order,
			}

			if hasFullAccess {
				lessonOut.VideoEmbedID = l.VideoEmbedID
				lessonOut.PrivateKey = l.PrivateKey

				if input.UserID > 0 {
					isCompleted, _ := u.progressReader.IsLessonCompleted(ctx, input.UserID, l.ID)
					lessonOut.IsCompleted = isCompleted
				}
			} else {
				lessonOut.VideoEmbedID = ""
				lessonOut.PrivateKey = nil
			}

			lessonsOut = append(lessonsOut, lessonOut)
		}

		modulesOut = append(modulesOut, ModuleOutput{
			ID:       m.ID,
			CourseID: m.CourseID,
			Title:    m.Title,
			Order:    m.Order,
			Lessons:  lessonsOut,
		})
	}

	progressPercent := 0
	if input.UserID > 0 && hasFullAccess {
		completed, total, err := u.progressReader.GetCourseProgress(ctx, input.UserID, course.ID)
		if err == nil && total > 0 {
			progressPercent = int(float64(completed) / float64(total) * 100)
		} else if err != nil {
			zap.L().Warn("Failed to calculate progress", zap.Int("userID", input.UserID), zap.Int("courseID", course.ID), zap.Error(err))
		}
	}

	course.IsPurchased = isPurchased

	zap.L().Info("GetFullCourse completed successfully",
		zap.Int("courseID", course.ID),
		zap.Bool("isPurchased", isPurchased),
		zap.Int("progressPercent", progressPercent))

	return &Output{
		Course:          course,
		Modules:         modulesOut,
		IsPurchased:     isPurchased,
		ProgressPercent: progressPercent,
	}, nil
}
