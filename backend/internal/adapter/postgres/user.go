package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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

func (r *UserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, full_name, role, created_at FROM users WHERE id = $1`,
		id,
	)

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

// UpdatePassword обновляет пароль пользователя
func (r *UserRepo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		passwordHash, userID,
	)
	return err
}

func (r *UserRepo) SearchByEmail(ctx context.Context, query string, limit int) ([]*domain.User, error) {
	log.Printf("[UserRepo] SearchByEmail called: query='%s', limit=%d", query, limit)

	// Исправлено: name -> full_name
	sql := `SELECT id, email, full_name, role, created_at FROM users WHERE email ILIKE $1 LIMIT $2`

	searchPattern := "%" + query + "%"
	log.Printf("[UserRepo] Executing SQL: %s with args: [%s, %d]", sql, searchPattern, limit)

	rows, err := r.db.QueryContext(ctx, sql, searchPattern, limit)
	if err != nil {
		log.Printf("[UserRepo] Query execution failed: %v", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	count := 0

	for rows.Next() {
		u := &domain.User{}
		var fullName string // Временная переменная для full_name

		// Исправлено: сканируем в fullName вместо u.Name
		err := rows.Scan(&u.ID, &u.Email, &fullName, &u.Role, &u.CreatedAt)
		if err != nil {
			log.Printf("[UserRepo] Scan failed: %v", err)
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		// Присваиваем значение в поле Name структуры User (если оно так называется в домене)
		u.Name = fullName

		users = append(users, u)
		count++
	}

	if err = rows.Err(); err != nil {
		log.Printf("[UserRepo] Rows iteration error: %v", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	log.Printf("[UserRepo] Search completed successfully. Found %d users.", count)
	return users, nil
}
