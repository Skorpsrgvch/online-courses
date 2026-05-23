package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	courseCreate "github.com/Skorpsrgvch/online-courses/internal/usecase/course/create"
	"go.uber.org/zap"
)

type CourseTxRepo struct {
	db *sql.DB
}

func NewCourseTxRepo(db *sql.DB) *CourseTxRepo {
	return &CourseTxRepo{db: db}
}

// SaveCourseWithModules создает курс и связанные с ним модули/уроки в одной транзакции.
func (r *CourseTxRepo) SaveCourseWithModules(ctx context.Context, course *domain.Course, modules []courseCreate.ModuleInput) error {
	zap.L().Info("Starting transaction: Create course with modules",
		zap.String("title", course.Title),
		zap.Int("modules_count", len(modules)),
	)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Откат при панике или ошибке (если коммит не прошел)
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	bonusesJSON, err := json.Marshal(course.Bonuses)
	if err != nil {
		_ = tx.Rollback()
		zap.L().Error("Failed to marshal bonuses JSON", zap.Error(err))
		return fmt.Errorf("failed to serialize bonuses: %w", err)
	}

	courseQuery := `
		INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, 
		                     contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, courseQuery,
		course.Title, course.Description, course.IsPublic, course.Price, course.AuthorID, course.IsActive,
		course.CoverImageURL, course.Contraindications, course.Recommendations, course.TargetAudience,
		course.CourseBasis, course.ClassBasis, bonusesJSON,
	).Scan(&course.ID)

	if err != nil {
		_ = tx.Rollback()
		zap.L().Error("Failed to insert course", zap.Error(err), zap.String("title", course.Title))
		return fmt.Errorf("failed to insert course: %w", err)
	}

	zap.L().Debug("Course inserted", zap.Int("course_id", course.ID))

	// Вставка модулей и уроков
	for i, mod := range modules {
		moduleQuery := `INSERT INTO modules (course_id, title, "order") VALUES ($1, $2, $3) RETURNING id`
		var moduleID int

		err = tx.QueryRowContext(ctx, moduleQuery, course.ID, mod.Title, mod.Order).Scan(&moduleID)
		if err != nil {
			_ = tx.Rollback()
			zap.L().Error("Failed to insert module", zap.Error(err), zap.String("title", mod.Title))
			return fmt.Errorf("failed to insert module %d: %w", i, err)
		}

		for j, lesson := range mod.Lessons {
			lessonQuery := `
				INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
				VALUES ($1, $2, $3, $4, $5, $6)
			`

			var pkVal interface{}
			if lesson.PrivateKey != nil && *lesson.PrivateKey != "" {
				pkVal = *lesson.PrivateKey
			} else {
				pkVal = nil
			}

			_, err = tx.ExecContext(ctx, lessonQuery, moduleID, lesson.Title, lesson.Description, lesson.VideoEmbedID, pkVal, lesson.Order)
			if err != nil {
				_ = tx.Rollback()
				zap.L().Error("Failed to insert lesson", zap.Error(err), zap.String("title", lesson.Title))
				return fmt.Errorf("failed to insert lesson %d in module %d: %w", j, i, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	zap.L().Info("Transaction committed successfully", zap.Int("course_id", course.ID))
	return nil
}
