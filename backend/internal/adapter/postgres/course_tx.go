package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	courseCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
)

type CourseTxRepo struct {
	db *sql.DB
}

func NewCourseTxRepo(db *sql.DB) *CourseTxRepo {
	return &CourseTxRepo{db: db}
}

func (r *CourseTxRepo) SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []courseCreate.ModuleInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сериализуем бонусы в JSON
	bonusesJSON, err := json.Marshal(course.Bonuses)
	if err != nil {
		return err
	}

	// 1. Создаём курс с новыми полями
	courseQuery := `
		INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, $9, $10, $11, $12, $13)
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
		course.Contraindications,
		course.Recommendations,
		course.TargetAudience,
		course.CourseBasis,
		course.ClassBasis,
		bonusesJSON,
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
			lessonQuery := `INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
							VALUES ($1, $2, $3, $4, $5, $6)`

			// Подготовка значения: если ключа нет, передаем nil
			var pkVal interface{}
			if lesson.PrivateKey != nil {
				pkVal = *lesson.PrivateKey
			} else {
				pkVal = nil
			}

			_, err = tx.ExecContext(ctx, lessonQuery,
				moduleID,
				lesson.Title,
				lesson.Description,
				lesson.VideoEmbedID,
				pkVal,
				lesson.Order,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
