package usercourses

import (
	"context"
	"errors"
	"math"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type Input struct {
	UserID int
}

type CourseProgress struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	IsPublic        bool   `json:"is_public"`
	CoverImageURL   string `json:"cover_image_url"`
	CompletedCount  int    `json:"completed_count"`
	TotalLessons    int    `json:"total_lessons"`
	ProgressPercent int    `json:"progress_percent"`
}

type Output struct {
	Courses []CourseProgress `json:"courses"`
}

type CourseLister interface {
	ListAll(ctx context.Context) ([]*domain.Course, error)
}

type PurchasesReader interface {
	GetUserCourseIDs(ctx context.Context, userID int) ([]int, error)
}

type ProgressReader interface {
	GetCompletedLessonsCount(ctx context.Context, userID, courseID int) (int, error)
	GetTotalLessonsCount(ctx context.Context, courseID int) (int, error)
}

type Usecase struct {
	courseLister    CourseLister
	purchasesReader PurchasesReader
	progressReader  ProgressReader
}

func NewUsecase(courseLister CourseLister, purchasesReader PurchasesReader, progressReader ProgressReader) (*Usecase, error) {
	if courseLister == nil || purchasesReader == nil || progressReader == nil {
		return nil, errors.New("all dependencies required")
	}
	return &Usecase{
		courseLister:    courseLister,
		purchasesReader: purchasesReader,
		progressReader:  progressReader,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	zap.L().Debug("GetUserCourses started", zap.Int("userID", input.UserID))

	courseIDs, err := u.purchasesReader.GetUserCourseIDs(ctx, input.UserID)
	if err != nil {
		zap.L().Error("Failed to get user course IDs", zap.Int("userID", input.UserID), zap.Error(err))
		return nil, err
	}

	if len(courseIDs) == 0 {
		zap.L().Info("No courses found for user", zap.Int("userID", input.UserID))
		return &Output{Courses: []CourseProgress{}}, nil
	}

	// Получаем все курсы
	allCourses, err := u.courseLister.ListAll(ctx)
	if err != nil {
		zap.L().Error("Failed to list all courses", zap.Error(err))
		return nil, err
	}

	// Фильтруем только купленные
	courseMap := make(map[int]*domain.Course)
	for _, c := range allCourses {
		for _, id := range courseIDs {
			if c.ID == id {
				courseMap[c.ID] = c
				break
			}
		}
	}

	// Считаем прогресс
	result := make([]CourseProgress, 0, len(courseIDs))
	for _, id := range courseIDs {
		c, ok := courseMap[id]
		if !ok {
			zap.L().Warn("Course ID from purchases not found in catalog", zap.Int("courseID", id))
			continue
		}

		total, errTotal := u.progressReader.GetTotalLessonsCount(ctx, c.ID)
		completed, errComp := u.progressReader.GetCompletedLessonsCount(ctx, input.UserID, c.ID)

		if errTotal != nil {
			zap.L().Warn("Failed to get total lessons count", zap.Int("courseID", c.ID), zap.Error(errTotal))
			total = 0
		}
		if errComp != nil {
			zap.L().Warn("Failed to get completed lessons count", zap.Int("courseID", c.ID), zap.Error(errComp))
			completed = 0
		}

		percent := 0
		if total > 0 {
			percent = int(math.Round(float64(completed) / float64(total) * 100))
		}

		result = append(result, CourseProgress{
			ID:              c.ID,
			Title:           c.Title,
			Description:     c.Description,
			Price:           c.Price,
			IsPublic:        c.IsPublic,
			CoverImageURL:   c.CoverImageURL,
			CompletedCount:  completed,
			TotalLessons:    total,
			ProgressPercent: percent,
		})
	}

	zap.L().Info("User courses retrieved successfully", zap.Int("userID", input.UserID), zap.Int("count", len(result)))
	return &Output{Courses: result}, nil
}
