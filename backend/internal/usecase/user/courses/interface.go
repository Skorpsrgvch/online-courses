package usercourses

import (
	"context"

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
