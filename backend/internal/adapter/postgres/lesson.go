package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type LessonRepo struct {
	db *sql.DB
}

func NewLessonRepo(db *sql.DB) *LessonRepo {
	return &LessonRepo{db: db}
}

func scanLesson(row *sql.Row) (*domain.Lesson, error) {
	var l domain.Lesson
	var pk sql.NullString

	err := row.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Description, &l.VideoEmbedID, &pk, &l.Order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLessonNotFound
		}
		return nil, fmt.Errorf("scan lesson: %w", err)
	}

	if pk.Valid {
		l.PrivateKey = &pk.String
	} else {
		l.PrivateKey = nil
	}

	return &l, nil
}

func (r *LessonRepo) GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, module_id, title, description, video_embed_id, private_key, "order"
		 FROM lessons WHERE module_id = $1 ORDER BY "order"`,
		moduleID,
	)
	if err != nil {
		return nil, fmt.Errorf("query lessons by module: %w", err)
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		var l domain.Lesson
		var pk sql.NullString

		if err := rows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Description, &l.VideoEmbedID, &pk, &l.Order); err != nil {
			return nil, fmt.Errorf("scan lesson row: %w", err)
		}

		if pk.Valid {
			l.PrivateKey = &pk.String
		} else {
			l.PrivateKey = nil
		}

		lessons = append(lessons, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lesson rows: %w", err)
	}

	zap.L().Debug("Lessons fetched by module", zap.Int("moduleID", moduleID), zap.Int("count", len(lessons)))
	return lessons, nil
}

func (r *LessonRepo) GetByID(ctx context.Context, id int) (*domain.Lesson, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, module_id, title, description, video_embed_id, private_key, "order"
		 FROM lessons WHERE id = $1`,
		id,
	)

	lesson, err := scanLesson(row)
	if err != nil {
		return nil, err
	}

	zap.L().Debug("Lesson fetched by ID", zap.Int("id", id))
	return lesson, nil
}

func (r *LessonRepo) Save(ctx context.Context, lesson *domain.Lesson) error {
	var pkVal interface{}
	if lesson.PrivateKey != nil {
		pkVal = *lesson.PrivateKey
	}

	query := `
		INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, lesson.ModuleID, lesson.Title, lesson.Description, lesson.VideoEmbedID, pkVal, lesson.Order).Scan(&lesson.ID)
	if err != nil {
		zap.L().Error("Failed to save lesson", zap.String("title", lesson.Title), zap.Error(err))
		return fmt.Errorf("save lesson: %w", err)
	}

	zap.L().Debug("Lesson saved", zap.Int("id", lesson.ID))
	return nil
}

func (r *LessonRepo) Update(ctx context.Context, lesson *domain.Lesson) error {
	var pkVal interface{}
	if lesson.PrivateKey != nil {
		pkVal = *lesson.PrivateKey
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE lessons SET title = $1, description = $2, video_embed_id = $3, private_key = $4, "order" = $5 WHERE id = $6`,
		lesson.Title, lesson.Description, lesson.VideoEmbedID, pkVal, lesson.Order, lesson.ID,
	)
	if err != nil {
		zap.L().Error("Failed to update lesson", zap.Int("id", lesson.ID), zap.Error(err))
		return fmt.Errorf("update lesson: %w", err)
	}

	zap.L().Debug("Lesson updated", zap.Int("id", lesson.ID))
	return nil
}

func (r *LessonRepo) UpdateOrderBatch(ctx context.Context, moduleID int, orders map[int]int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	for id, order := range orders {
		select {
		case <-ctx.Done():
			_ = tx.Rollback()
			return ctx.Err()
		default:
		}

		_, err := tx.ExecContext(ctx, `UPDATE lessons SET "order" = $1 WHERE id = $2 AND module_id = $3`, order, id, moduleID)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("update order for lesson %d: %w", id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	zap.L().Debug("Lesson orders updated batch", zap.Int("moduleID", moduleID), zap.Int("count", len(orders)))
	return nil
}

func (r *LessonRepo) Delete(ctx context.Context, lessonID int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM lessons WHERE id = $1`, lessonID)
	if err != nil {
		zap.L().Error("Failed to delete lesson", zap.Int("id", lessonID), zap.Error(err))
		return fmt.Errorf("delete lesson: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		zap.L().Warn("Attempted to delete non-existent lesson", zap.Int("id", lessonID))
		return domain.ErrLessonNotFound
	}

	zap.L().Debug("Lesson deleted", zap.Int("id", lessonID))
	return nil
}
