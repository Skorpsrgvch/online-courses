package getfull

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	CourseID int
	UserID   int    // 0 — если неавторизован
	Role     string // "admin", "user"
}

type LessonOutput struct {
	ID           int     `json:"id"`
	ModuleID     int     `json:"module_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key"`
	Order        int     `json:"order"`
}

type ModuleOutput struct {
	ID       int            `json:"id"`
	CourseID int            `json:"course_id"`
	Title    string         `json:"title"`
	Order    int            `json:"order"`
	Lessons  []LessonOutput `json:"lessons"`
}

type Output struct {
	Course      *domain.Course `json:"course"`
	Modules     []ModuleOutput `json:"modules"`
	IsPurchased bool           `json:"is_purchased"`
}

type Usecase struct {
	courseReader    CourseReader
	moduleReader    ModuleReader
	lessonReader    LessonReader
	purchaseChecker PurchaseChecker
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
