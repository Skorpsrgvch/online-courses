package update

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	ID                int
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
}

type Usecase struct {
	courseUpdater CourseUpdater
	courseFinder  CourseFinder
}

type CourseUpdater interface {
	Update(ctx context.Context, course *domain.Course) error
}

type CourseFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}
