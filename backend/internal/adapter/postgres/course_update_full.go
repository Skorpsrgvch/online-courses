package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/updatefullcourse"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type CourseFullRepo struct {
	db *sql.DB
}

func NewCourseFullRepo(db *sql.DB) *CourseFullRepo {
	return &CourseFullRepo{db: db}
}

func (r *CourseFullRepo) UpdateFullWithModules(ctx context.Context, input updatefullcourse.Input) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Failed to begin transaction for full course update", zap.Int("courseID", input.CourseID), zap.Error(err))
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	bonusesJSON, err := json.Marshal(input.Bonuses)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("marshal bonuses: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE courses SET 
			title = $1, description = $2, is_public = $3, price = $4, is_active = $5,
			cover_image_url = NULLIF($6, ''), contraindications = $7, recommendations = $8,
			target_audience = $9, course_basis = $10, class_basis = $11, bonuses = $12
		WHERE id = $13`,
		input.Title, input.Description, input.IsPublic, input.Price, input.IsActive,
		input.CoverImageURL, input.Contraindications, input.Recommendations,
		input.TargetAudience, input.CourseBasis, input.ClassBasis, bonusesJSON, input.CourseID,
	)
	if err != nil {
		_ = tx.Rollback()
		zap.L().Error("Failed to update course details", zap.Int("courseID", input.CourseID), zap.Error(err))
		return fmt.Errorf("update course: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT id FROM modules WHERE course_id = $1`, input.CourseID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("query existing modules: %w", err)
	}
	existingModuleIDs := make(map[int]bool)
	var dbModID int
	for rows.Next() {
		if err := rows.Scan(&dbModID); err != nil {
			rows.Close()
			_ = tx.Rollback()
			return fmt.Errorf("scan module: %w", err)
		}
		existingModuleIDs[dbModID] = true
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		_ = tx.Rollback()
		return err
	}
	rows.Close()

	for _, mod := range input.Modules {
		select {
		case <-ctx.Done():
			_ = tx.Rollback()
			return ctx.Err()
		default:
		}

		var modID int
		if mod.ID == 0 {
			err = tx.QueryRowContext(ctx, `
				INSERT INTO modules (course_id, title, "order") VALUES ($1, $2, $3) RETURNING id`,
				input.CourseID, mod.Title, mod.Order).Scan(&modID)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("insert module: %w", err)
			}
		} else {
			modID = mod.ID
			_, err = tx.ExecContext(ctx, `
				UPDATE modules SET title = $1, "order" = $2 WHERE id = $3`,
				mod.Title, mod.Order, modID)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("update module: %w", err)
			}
			delete(existingModuleIDs, modID)
		}

		if err := r.syncLessons(ctx, tx, modID, mod.Lessons); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("sync lessons for module %d: %w", modID, err)
		}
	}

	for id := range existingModuleIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM modules WHERE id = $1`, id); err != nil {
			_ = tx.Rollback()
			zap.L().Error("Failed to delete orphan module", zap.Int("moduleID", id), zap.Error(err))
			return fmt.Errorf("delete module: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Int("courseID", input.CourseID), zap.Error(err))
		return fmt.Errorf("commit transaction: %w", err)
	}

	zap.L().Info("Course fully updated with modules and lessons", zap.Int("courseID", input.CourseID))
	return nil
}

func (r *CourseFullRepo) syncLessons(ctx context.Context, tx *sql.Tx, moduleID int, lessons []updatefullcourse.LessonInput) error {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM lessons WHERE module_id = $1`, moduleID)
	if err != nil {
		return fmt.Errorf("query existing lessons: %w", err)
	}
	existingLessonIDs := make(map[int]bool)
	var dbLessID int
	for rows.Next() {
		if err := rows.Scan(&dbLessID); err != nil {
			rows.Close()
			return fmt.Errorf("scan lesson: %w", err)
		}
		existingLessonIDs[dbLessID] = true
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	for _, less := range lessons {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var pkVal interface{}
		if less.PrivateKey != nil {
			pkVal = *less.PrivateKey
		}

		if less.ID == 0 {
			err = tx.QueryRowContext(ctx, `
				INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
				moduleID, less.Title, less.Description, less.VideoEmbedID, pkVal, less.Order).Scan(&less.ID)
			if err != nil {
				return fmt.Errorf("insert lesson: %w", err)
			}
		} else {
			_, err = tx.ExecContext(ctx, `
				UPDATE lessons SET 
					title = $1, description = $2, video_embed_id = $3, private_key = $4, "order" = $5 
				WHERE id = $6`,
				less.Title, less.Description, less.VideoEmbedID, pkVal, less.Order, less.ID)
			if err != nil {
				return fmt.Errorf("update lesson: %w", err)
			}
			delete(existingLessonIDs, less.ID)
		}
	}

	for id := range existingLessonIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM lessons WHERE id = $1`, id); err != nil {
			return fmt.Errorf("delete lesson: %w", err)
		}
	}

	return nil
}

func (r *CourseFullRepo) ReorderModules(ctx context.Context, courseID int, orderMap map[int]int) error {
	if len(orderMap) == 0 {
		return nil
	}

	var ids []int
	var args []interface{}
	for id, order := range orderMap {
		ids = append(ids, id)
		args = append(args, id, order)
	}

	caseStmt := "CASE id "
	for i := 0; i < len(ids); i++ {
		caseStmt += fmt.Sprintf("WHEN $%d THEN $%d ", i*2+1, i*2+2)
	}
	caseStmt += "END"

	query := fmt.Sprintf(`UPDATE modules SET "order" = %s WHERE course_id = $%d AND id = ANY($%d::int[])`,
		caseStmt, len(args)+1, len(args)+2)

	args = append(args, courseID, pq.Array(ids))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		zap.L().Error("Failed to reorder modules", zap.Int("courseID", courseID), zap.Error(err))
		return fmt.Errorf("reorder modules: %w", err)
	}

	zap.L().Debug("Modules reordered", zap.Int("courseID", courseID), zap.Int("count", len(orderMap)))
	return nil
}

func (r *CourseFullRepo) ReorderLessons(ctx context.Context, moduleID int, orderMap map[int]int) error {
	if len(orderMap) == 0 {
		return nil
	}

	var ids []int
	var args []interface{}
	for id, order := range orderMap {
		ids = append(ids, id)
		args = append(args, id, order)
	}

	caseStmt := "CASE id "
	for i := 0; i < len(ids); i++ {
		caseStmt += fmt.Sprintf("WHEN $%d THEN $%d ", i*2+1, i*2+2)
	}
	caseStmt += "END"

	query := fmt.Sprintf(`UPDATE lessons SET "order" = %s WHERE module_id = $%d AND id = ANY($%d::int[])`,
		caseStmt, len(args)+1, len(args)+2)

	args = append(args, moduleID, pq.Array(ids))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		zap.L().Error("Failed to reorder lessons", zap.Int("moduleID", moduleID), zap.Error(err))
		return fmt.Errorf("reorder lessons: %w", err)
	}

	zap.L().Debug("Lessons reordered", zap.Int("moduleID", moduleID), zap.Int("count", len(orderMap)))
	return nil
}
