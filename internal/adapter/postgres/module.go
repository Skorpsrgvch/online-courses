package postgres

import (
	"context"
	"database/sql"

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

func (r *ModuleRepo) Save(ctx context.Context, module *domain.Module) error {
	query := `
		INSERT INTO modules (course_id, title, "order")
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, module.CourseID, module.Title, module.Order).Scan(&module.ID)
}
