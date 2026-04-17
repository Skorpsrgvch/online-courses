package usercourses

import (
	"context"
	"errors"
	"math"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
	courseIDs, err := u.purchasesReader.GetUserCourseIDs(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if len(courseIDs) == 0 {
		return &Output{Courses: []CourseProgress{}}, nil
	}

	// Получаем все курсы
	allCourses, err := u.courseLister.ListAll(ctx)
	if err != nil {
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
	var result []CourseProgress
	for _, id := range courseIDs {
		c, ok := courseMap[id]
		if !ok {
			continue
		}

		total, _ := u.progressReader.GetTotalLessonsCount(ctx, c.ID)
		completed, _ := u.progressReader.GetCompletedLessonsCount(ctx, input.UserID, c.ID)

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

	return &Output{Courses: result}, nil
}
