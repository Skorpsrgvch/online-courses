package updatefullcourse

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type LessonInput struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	VideoEmbedID string  `json:"video_embed_id"`
	PrivateKey   *string `json:"private_key,omitempty"`
	Order        int     `json:"order"`
}

type ModuleInput struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Order   int           `json:"order"`
	Lessons []LessonInput `json:"lessons"`
}

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

type CourseRepository interface {
	UpdateFullWithModules(ctx context.Context, input Input) error
	ReorderModules(ctx context.Context, courseID int, orderMap map[int]int) error
	ReorderLessons(ctx context.Context, courseID int, orderMap map[int]int) error
}
