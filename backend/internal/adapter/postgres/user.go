package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User, passwordHash string) error {
	if user == nil || user.Email == "" || passwordHash == "" {
		return errors.New("invalid user data or password hash")
	}

	query := `
		INSERT INTO users (email, full_name, role, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	err := r.db.QueryRowContext(ctx, query,
		user.Email, user.Name, user.Role, passwordHash, user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return common.BadRequestError(fmt.Sprintf("пользователь с email '%s' уже существует", user.Email))
		}
		zap.L().Error("UserRepo.CreateUser failed", zap.String("email", user.Email), zap.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	zap.L().Info("User created", zap.Int("id", user.ID), zap.String("email", user.Email))
	return nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1 LIMIT 1`, email).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *UserRepo) GetUserByEmailForCheck(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, full_name, role, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	u.CreatedAt = createdAt.Time
	return &u, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}

	query := `SELECT id, email, full_name, role, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	u.CreatedAt = createdAt.Time
	return &u, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	query := `SELECT id, email, full_name, role, created_at FROM users WHERE email = $1 AND password_hash = $2`
	row := r.db.QueryRowContext(ctx, query, email, passwordHash)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	u.CreatedAt = createdAt.Time
	return &u, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	query := `SELECT id, email, full_name, role, password_hash, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	var passwordHash string
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &passwordHash, &createdAt)

	if err == sql.ErrNoRows {
		return nil, "", domain.ErrUserNotFound
	}
	if err != nil {
		return nil, "", err
	}

	u.CreatedAt = createdAt.Time
	return &u, passwordHash, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	if userID <= 0 || passwordHash == "" {
		return errors.New("invalid user ID or password hash")
	}

	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) UpdateName(ctx context.Context, userID int, name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := r.db.ExecContext(ctx, `UPDATE users SET full_name = $1 WHERE id = $2`, name, userID)
	return err
}

func (r *UserRepo) UpdateEmail(ctx context.Context, userID int, email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	query := `UPDATE users SET email = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, email, userID)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("email '%s' already exists", email)
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetPasswordHash(ctx context.Context, userID int) (string, error) {
	var hash string
	err := r.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)
	if err == sql.ErrNoRows {
		return "", domain.ErrUserNotFound
	}
	return hash, err
}

func (r *UserRepo) SearchByEmail(ctx context.Context, query string, limit int) ([]*domain.User, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	if limit <= 0 {
		limit = 10
	}

	sqlQuery := `SELECT id, email, full_name, role, created_at FROM users WHERE email ILIKE $1 LIMIT $2`
	searchPattern := "%" + query + "%"

	rows, err := r.db.QueryContext(ctx, sqlQuery, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var fullName string
		var createdAt sql.NullTime

		if err := rows.Scan(&u.ID, &u.Email, &fullName, &u.Role, &createdAt); err != nil {
			return nil, err
		}

		u.Name = fullName
		if createdAt.Valid {
			u.CreatedAt = createdAt.Time
		}
		users = append(users, &u)
	}

	return users, rows.Err()
}
