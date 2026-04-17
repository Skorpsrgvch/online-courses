// internal/adapter/postgres/course_tx.go
package postgres

import (
	"context"
	"database/sql"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	courseCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
)

// CourseTxRepo — репозиторий курсов с поддержкой транзакций
type CourseTxRepo struct {
	db *sql.DB
}

func NewCourseTxRepo(db *sql.DB) *CourseTxRepo {
	return &CourseTxRepo{db: db}
}

// SaveCourseWithModules создаёт курс, модули и уроки в одной транзакции
func (r *CourseTxRepo) SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []courseCreate.ModuleInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Создаём курс
	courseQuery := `
		INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''))
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, courseQuery,
		course.Title,
		course.Description,
		course.IsPublic,
		course.Price,
		course.AuthorID,
		course.IsActive,
		course.CoverImageURL,
	).Scan(&course.ID)
	if err != nil {
		return err
	}

	// 2. Создаём модули и уроки
	for _, mod := range modules {
		moduleQuery := `
			INSERT INTO modules (course_id, title, "order")
			VALUES ($1, $2, $3)
			RETURNING id
		`
		var moduleID int
		err = tx.QueryRowContext(ctx, moduleQuery, course.ID, mod.Title, mod.Order).Scan(&moduleID)
		if err != nil {
			return err
		}

		// 3. Создаём уроки для каждого модуля
		for _, lesson := range mod.Lessons {
			lessonQuery := `
				INSERT INTO lessons (module_id, title, description, lesson_type, video_embed_id, article_content, "order")
				VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7)
			`
			_, err = tx.ExecContext(ctx, lessonQuery,
				moduleID,
				lesson.Title,
				lesson.Description,
				string(lesson.LessonType),
				lesson.VideoEmbedID,
				lesson.ArticleContent,
				lesson.Order,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
