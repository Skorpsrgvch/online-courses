package update

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

func NewUsecase(updater CourseUpdater, finder CourseFinder) (*Usecase, error) {
	if updater == nil || finder == nil {
		return nil, errors.New("dependencies required")
	}
	return &Usecase{
		courseUpdater: updater,
		courseFinder:  finder,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	// 1. Проверяем существование курса (чтобы не потерять данные, если какие-то поля не переданы)
	existing, err := u.courseFinder.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	// 2. Логика для картинки: если новую не передали, оставляем старую
	coverImageURL := input.CoverImageURL
	if coverImageURL == "" {
		coverImageURL = existing.CoverImageURL
	}

	// 3. Создаем объект курса с обновленными данными.
	// ВАЖНО: RestoreCourse возвращает только *Course, без ошибки.
	updated := domain.RestoreCourse(
		input.ID,
		input.Title,
		input.Description,
		input.IsPublic,
		input.Price,
		input.AuthorID,
		input.IsActive,
		coverImageURL,
		input.Contraindications,
		input.Recommendations,
		input.TargetAudience,
		input.CourseBasis,
		input.ClassBasis,
	)

	// 4. Обработка бонусов
	if input.Bonuses == nil {
		// Если бонусы не передали, оставляем пустой срез, а не nil (для корректного JSON)
		updated.Bonuses = []domain.BonusItem{}
	} else {
		updated.Bonuses = input.Bonuses
	}

	// 5. Сохраняем изменения
	return u.courseUpdater.Update(ctx, updated)
}
