package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ServiceRepo struct {
	db *sql.DB
}

func NewServiceRepo(db *sql.DB) *ServiceRepo {
	return &ServiceRepo{db: db}
}

func (r *ServiceRepo) GetAll(ctx context.Context) ([]*domain.Service, error) {
	query := `
		SELECT id, title, description, price, duration_minutes, created_at, updated_at
		FROM services
		ORDER BY price ASC, id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*domain.Service
	for rows.Next() {
		var s domain.Service
		err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Price, &s.Duration, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, &s)
	}
	return services, rows.Err()
}

func (r *ServiceRepo) GetByID(ctx context.Context, id int) (*domain.Service, error) {
	query := `
		SELECT id, title, description, price, duration_minutes, created_at, updated_at
		FROM services WHERE id = $1
	`
	var s domain.Service
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Title, &s.Description, &s.Price, &s.Duration, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrServiceNotFound
	}
	return &s, err
}

func (r *ServiceRepo) Create(ctx context.Context, s *domain.Service) error {
	query := `
		INSERT INTO services (title, description, price, duration_minutes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	now := time.Now().UTC()
	return r.db.QueryRowContext(ctx, query, s.Title, s.Description, s.Price, s.Duration, now, now).Scan(&s.ID)
}

func (r *ServiceRepo) Update(ctx context.Context, s *domain.Service) error {
	query := `
		UPDATE services
		SET title = $1, description = $2, price = $3, duration_minutes = $4, updated_at = $5
		WHERE id = $6
	`
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, s.Title, s.Description, s.Price, s.Duration, now, s.ID)
	return err
}

func (r *ServiceRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM services WHERE id = $1`, id)
	return err
}
