package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/course/updatefullcourse"
	"github.com/lib/pq"
)

type CourseFullRepo struct {
	db *sql.DB
}

func NewCourseFullRepo(db *sql.DB) *CourseFullRepo {
	return &CourseFullRepo{db: db}
}

// UpdateFullWithModules выполняет полное обновление курса, модулей и уроков в транзакции
func (r *CourseFullRepo) UpdateFullWithModules(ctx context.Context, input updatefullcourse.Input) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Обновляем сам курс
	bonusesJSON, err := json.Marshal(input.Bonuses)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE courses SET 
			title = $1, 
			description = $2,
			is_public = $3,
			price = $4,
			is_active = $5,
			cover_image_url = NULLIF($6, ''),
			contraindications = $7,
			recommendations = $8,
			target_audience = $9,
			course_basis = $10,
			class_basis = $11,
			bonuses = $12
		WHERE id = $13`,
		input.Title,
		input.Description,
		input.IsPublic,
		input.Price,
		input.IsActive,
		input.CoverImageURL,
		input.Contraindications,
		input.Recommendations,
		input.TargetAudience,
		input.CourseBasis,
		input.ClassBasis,
		bonusesJSON,
		input.CourseID,
	)
	if err != nil {
		return err
	}

	// 2. Получаем текущие модули для выявления удаленных
	rows, err := tx.QueryContext(ctx, `SELECT id FROM modules WHERE course_id = $1`, input.CourseID)
	if err != nil {
		return err
	}
	existingModuleIDs := make(map[int]bool)
	var dbModID int
	for rows.Next() {
		if err := rows.Scan(&dbModID); err != nil {
			rows.Close()
			return err
		}
		existingModuleIDs[dbModID] = true
	}
	rows.Close()

	// 3. Обрабатываем модули из запроса
	for _, mod := range input.Modules {
		var modID int

		if mod.ID == 0 {
			// Создаем новый модуль
			err = tx.QueryRowContext(ctx, `
				INSERT INTO modules (course_id, title, "order") VALUES ($1, $2, $3) RETURNING id`,
				input.CourseID, mod.Title, mod.Order).Scan(&modID)
			if err != nil {
				return err
			}
		} else {
			// Обновляем существующий
			modID = mod.ID
			_, err = tx.ExecContext(ctx, `
				UPDATE modules SET title = $1, "order" = $2 WHERE id = $3`,
				mod.Title, mod.Order, modID)
			if err != nil {
				return err
			}
			// Помечаем как обработанный
			delete(existingModuleIDs, modID)
		}

		// 4. Синхронизируем уроки
		if err := r.syncLessons(ctx, tx, modID, mod.Lessons); err != nil {
			return err
		}
	}

	// 5. Удаляем лишние модули (каскад удалит уроки)
	for id := range existingModuleIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM modules WHERE id = $1`, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *CourseFullRepo) syncLessons(ctx context.Context, tx *sql.Tx, moduleID int, lessons []updatefullcourse.LessonInput) error {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM lessons WHERE module_id = $1`, moduleID)
	if err != nil {
		return err
	}
	existingLessonIDs := make(map[int]bool)
	var dbLessID int
	for rows.Next() {
		if err := rows.Scan(&dbLessID); err != nil {
			rows.Close()
			return err
		}
		existingLessonIDs[dbLessID] = true
	}
	rows.Close()

	for _, less := range lessons {
		var lessID int
		var pkVal interface{}
		if less.PrivateKey != nil {
			pkVal = *less.PrivateKey
		}

		if less.ID == 0 {
			err = tx.QueryRowContext(ctx, `
				INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
				moduleID, less.Title, less.Description, less.VideoEmbedID, pkVal, less.Order).Scan(&lessID)
			if err != nil {
				return err
			}
		} else {
			lessID = less.ID
			_, err = tx.ExecContext(ctx, `
				UPDATE lessons SET 
					title = $1, 
					description = $2, 
					video_embed_id = $3, 
					private_key = $4, 
					"order" = $5 
				WHERE id = $6`,
				less.Title, less.Description, less.VideoEmbedID, pkVal, less.Order, lessID)
			if err != nil {
				return err
			}
			delete(existingLessonIDs, lessID)
		}
	}

	for id := range existingLessonIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM lessons WHERE id = $1`, id); err != nil {
			return err
		}
	}

	return nil
}

func (r *CourseFullRepo) ReorderModules(ctx context.Context, courseID int, orderMap map[int]int) error {
	if len(orderMap) == 0 {
		return nil
	}

	// Формируем SQL: UPDATE modules SET "order" = CASE id WHEN $1 THEN $2 WHEN $3 THEN $4 ... END WHERE id IN (...)
	var ids []int
	var args []interface{}
	argIndex := 1

	for id, order := range orderMap {
		ids = append(ids, id)
		args = append(args, id, order)
		argIndex += 2
	}

	// Создаем строку условий CASE
	caseStmt := "CASE id "
	for i := 0; i < len(ids); i++ {
		caseStmt += fmt.Sprintf("WHEN $%d THEN $%d ", i*2+1, i*2+2)
	}
	caseStmt += "END"

	query := fmt.Sprintf(`UPDATE modules SET "order" = %s WHERE course_id = $%d AND id = ANY($%d::int[])`,
		caseStmt,
		len(args)+1,
		len(args)+2,
	)

	args = append(args, courseID, pq.Array(ids))

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// ReorderLessons обновляет порядок уроков внутри конкретного модуля
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
		caseStmt,
		len(args)+1,
		len(args)+2,
	)

	args = append(args, moduleID, pq.Array(ids))

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
