package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, email, full_name, role, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	u.CreatedAt = createdAt.Time
	return &u, nil
}

// GetUserByEmailAndPassword проверяет учётные данные
func (r *UserRepo) GetUserByEmailAndPassword(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	query := `SELECT id, email, full_name, role, created_at FROM users WHERE email = $1 AND password_hash = $2`
	row := r.db.QueryRowContext(ctx, query, email, passwordHash)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	u.CreatedAt = createdAt.Time
	return &u, nil
}

// CreateUser создаёт нового пользователя
func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User, passwordHash string) error {
	query := `
		INSERT INTO users (email, full_name, role, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		user.Email, user.Name, user.Role, passwordHash, user.CreatedAt,
	).Scan(&user.ID)
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	query := `SELECT id, email, full_name, role, password_hash, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	var passwordHash string
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &passwordHash, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", domain.ErrUserNotFound
		}
		return nil, "", err
	}
	u.CreatedAt = createdAt.Time
	return &u, passwordHash, nil
}

func (r *CourseRepo) Update(ctx context.Context, course *domain.Course) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE courses SET title = $1, description = $2, is_public = $3, price = $4, is_active = $5 WHERE id = $6`,
		course.Title, course.Description, course.IsPublic, course.Price, course.IsActive, course.ID,
	)
	return err
}

// SetInactive деактивирует курс (логическое удаление)
func (r *CourseRepo) SetInactive(ctx context.Context, courseID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE courses SET is_active = false WHERE id = $1`,
		courseID,
	)
	return err
}
