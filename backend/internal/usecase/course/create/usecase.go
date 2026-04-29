package create

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

func NewUsecase(courseSaver CourseSaver) (*Usecase, error) {
	if courseSaver == nil {
		return nil, errors.New("courseSaver is required")
	}
	return &Usecase{courseSaver: courseSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	// Защита от nil для бонусов
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
		return err
	}

	course.Bonuses = bonuses

	return u.courseSaver.Save(ctx, course)
}
