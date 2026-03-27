package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
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
		return nil, err
	}
	defer rows.Close()

	var modules []*domain.Module
	for rows.Next() {
		var m domain.Module
		err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Order)
		if err != nil {
			return nil, err
		}
		modules = append(modules, &m)
	}
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
		return nil, err
	}
	return &m, nil
}

func (r *ModuleRepo) Save(ctx context.Context, module *domain.Module) error {
	query := `
		INSERT INTO modules (course_id, title, "order")
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, module.CourseID, module.Title, module.Order).Scan(&module.ID)
}

// Update обновляет модуль
func (r *ModuleRepo) Update(ctx context.Context, module *domain.Module) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE modules 
		 SET title = $1, "order" = $2 
		 WHERE id = $3`,
		module.Title,
		module.Order,
		module.ID,
	)
	return err
}

// Delete удаляет модуль и все его уроки (каскадно)
func (r *ModuleRepo) Delete(ctx context.Context, moduleID int) error {
	// Сначала удаляем уроки (если нет ON DELETE CASCADE)
	_, err := r.db.ExecContext(ctx, `DELETE FROM lessons WHERE module_id = $1`, moduleID)
	if err != nil {
		return err
	}
	// Теперь удаляем модуль
	_, err = r.db.ExecContext(ctx, `DELETE FROM modules WHERE id = $1`, moduleID)
	return err
}
