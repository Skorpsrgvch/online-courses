package createwithmodules

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
)

func NewUsecase(courseSaver create.CourseModuleSaver) (*Usecase, error) {
	if courseSaver == nil {
		return nil, errors.New("courseSaver is required")
	}
	return &Usecase{courseSaver: courseSaver}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	// 1. Создаем базовый объект курса, передавая ВСЕ текстовые поля в конструктор
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

	// 2. Защищаем бонусы от nil
	if input.Bonuses == nil {
		course.Bonuses = []domain.BonusItem{}
	} else {
		course.Bonuses = input.Bonuses
	}

	// 3. Сохраняем курс вместе с модулями
	// Логика сохранения модулей находится внутри репозитория (SaveCourseWithModules)
	return u.courseSaver.SaveCourseWithModules(ctx, course, input.Modules)
}
