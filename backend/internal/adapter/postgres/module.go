package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type ModuleRepo struct {
	db *sql.DB
}

func NewModuleRepo(db *sql.DB) *ModuleRepo {
	return &ModuleRepo{db: db}
}

func (r *ModuleRepo) GetByCourseID(ctx context.Context, courseID int) ([]*domain.Module, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, course_id, title, "order" FROM modules WHERE course_id = $1 ORDER BY "order"`,
		courseID,
	)
	if err != nil {
		return nil, fmt.Errorf("query modules by course: %w", err)
	}
	defer rows.Close()

	var modules []*domain.Module
	for rows.Next() {
		var m domain.Module
		if err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Order); err != nil {
			return nil, fmt.Errorf("scan module row: %w", err)
		}
		modules = append(modules, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate module rows: %w", err)
	}

	zap.L().Debug("Modules fetched by course", zap.Int("courseID", courseID), zap.Int("count", len(modules)))
	return modules, nil
}

func (r *ModuleRepo) GetByID(ctx context.Context, id int) (*domain.Module, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, course_id, title, "order" FROM modules WHERE id = $1`,
		id,
	)

	var m domain.Module
	err := row.Scan(&m.ID, &m.CourseID, &m.Title, &m.Order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrModuleNotFound
		}
		return nil, fmt.Errorf("scan module: %w", err)
	}

	zap.L().Debug("Module fetched by ID", zap.Int("id", id))
	return &m, nil
}

func (r *ModuleRepo) Save(ctx context.Context, module *domain.Module) error {
	query := `INSERT INTO modules (course_id, title, "order") VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, module.CourseID, module.Title, module.Order).Scan(&module.ID)
	if err != nil {
		zap.L().Error("Failed to save module", zap.String("title", module.Title), zap.Error(err))
		return fmt.Errorf("save module: %w", err)
	}

	zap.L().Debug("Module saved", zap.Int("id", module.ID))
	return nil
}

func (r *ModuleRepo) Update(ctx context.Context, module *domain.Module) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE modules SET title = $1, "order" = $2 WHERE id = $3`,
		module.Title, module.Order, module.ID,
	)
	if err != nil {
		zap.L().Error("Failed to update module", zap.Int("id", module.ID), zap.Error(err))
		return fmt.Errorf("update module: %w", err)
	}

	zap.L().Debug("Module updated", zap.Int("id", module.ID))
	return nil
}

func (r *ModuleRepo) UpdateOrderBatch(ctx context.Context, courseID int, orders map[int]int) error {
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

		_, err := tx.ExecContext(ctx, `UPDATE modules SET "order" = $1 WHERE id = $2 AND course_id = $3`, order, id, courseID)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("update order for module %d: %w", id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	zap.L().Debug("Module orders updated batch", zap.Int("courseID", courseID), zap.Int("count", len(orders)))
	return nil
}

func (r *ModuleRepo) Delete(ctx context.Context, moduleID int) error {
	// Каскадное удаление уроков происходит на уровне БД благодаря ON DELETE CASCADE
	_, err := r.db.ExecContext(ctx, `DELETE FROM modules WHERE id = $1`, moduleID)
	if err != nil {
		zap.L().Error("Failed to delete module", zap.Int("id", moduleID), zap.Error(err))
		return fmt.Errorf("delete module: %w", err)
	}

	zap.L().Debug("Module deleted", zap.Int("id", moduleID))
	return nil
}
