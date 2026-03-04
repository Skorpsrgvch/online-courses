// internal/adapter/postgres/course.go
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type CourseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *CourseRepo {
	return &CourseRepo{db: db}
}

// GetByID возвращает курс по ID
func (r *CourseRepo) GetByID(ctx context.Context, id int) (*domain.Course, error) {
	query := `
		SELECT id, title, description, is_public, price, author_id, is_active
		FROM courses
		WHERE id = $1 AND is_active = true
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var c domain.Course
	err := row.Scan(
		&c.ID,
		&c.Title,
		&c.Description,
		&c.IsPublic,
		&c.Price,
		&c.AuthorID,
		&c.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCourseNotFound
		}
		return nil, err
	}
	return &c, nil
}

// Save создаёт новый курс
func (r *CourseRepo) Save(ctx context.Context, course *domain.Course) error {
	query := `
		INSERT INTO courses (title, description, is_public, price, author_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		course.Title,
		course.Description,
		course.IsPublic,
		course.Price,
		course.AuthorID,
		course.IsActive,
	).Scan(&course.ID)
}

// Update обновляет курс (все поля, кроме ID)
func (r *CourseRepo) Update(ctx context.Context, course *domain.Course) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE courses 
		 SET title = $1, description = $2, is_public = $3, price = $4, is_active = $5 
		 WHERE id = $6`,
		course.Title,
		course.Description,
		course.IsPublic,
		course.Price,
		course.IsActive,
		course.ID,
	)
	return err
}

// SetInactive выполняет логическое удаление курса
func (r *CourseRepo) SetInactive(ctx context.Context, courseID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE courses SET is_active = false WHERE id = $1`,
		courseID,
	)
	return err
}
