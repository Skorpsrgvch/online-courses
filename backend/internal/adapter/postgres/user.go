package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/adapter/http/common"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	return &UserRepo{db: db}
}

// CreateUser создаёт нового пользователя
func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User, passwordHash string) error {
	log.Printf("[UserRepo] CreateUser called: email='%s', name='%s', role='%s'", user.Email, user.Name, user.Role)

	if user == nil {
		return errors.New("user object is nil")
	}
	if user.Email == "" {
		return errors.New("email cannot be empty")
	}
	if passwordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	query := `
		INSERT INTO users (email, full_name, role, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	err := r.db.QueryRowContext(ctx, query,
		user.Email, user.Name, user.Role, passwordHash, user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		log.Printf("[UserRepo] CreateUser error: %v", err)

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
			strings.Contains(err.Error(), "users_email_key") {
			return common.BadRequestError(fmt.Sprintf("пользователь с '%s' уже существует", user.Email))
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("[UserRepo] CreateUser success: created user with ID=%d", user.ID)
	return nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	log.Printf("[UserRepo] ExistsByEmail called: email='%s'", email)

	var id int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1 LIMIT 1`, email).Scan(&id)

	if err == sql.ErrNoRows {
		log.Printf("[UserRepo] ExistsByEmail result: false (email '%s' not found)", email)
		return false, nil
	}
	if err != nil {
		log.Printf("[UserRepo] ExistsByEmail error: %v", err)
		return false, fmt.Errorf("database query failed: %w", err)
	}

	log.Printf("[UserRepo] ExistsByEmail result: true (user ID=%d found)", id)
	return true, nil
}

func (r *UserRepo) GetUserByEmailForCheck(ctx context.Context, email string) (*domain.User, error) {
	log.Printf("[UserRepo] GetUserByEmailForCheck called: email='%s'", email)

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	query := `SELECT id, email, full_name, role, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[UserRepo] GetUserByEmailForCheck result: user not found (email is free)")
			return nil, domain.ErrUserNotFound
		}
		log.Printf("[UserRepo] GetUserByEmailForCheck error: %v", err)
		return nil, fmt.Errorf("database scan failed: %w", err)
	}

	u.CreatedAt = createdAt.Time
	log.Printf("[UserRepo] GetUserByEmailForCheck success: found user ID=%d (email is busy)", u.ID)
	return &u, nil
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	log.Printf("[UserRepo] GetUserByID called: id=%d", id)

	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}

	query := `SELECT id, email, full_name, role, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[UserRepo] GetUserByID result: user not found (id=%d)", id)
			return nil, domain.ErrUserNotFound
		}
		log.Printf("[UserRepo] GetUserByID error: %v", err)
		return nil, fmt.Errorf("database scan failed: %w", err)
	}

	u.CreatedAt = createdAt.Time
	log.Printf("[UserRepo] GetUserByID success: found user '%s' (role=%s)", u.Email, u.Role)
	return &u, nil
}

// GetUserByEmailAndPassword проверяет учётные данные
func (r *UserRepo) GetUserByEmailAndPassword(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	log.Printf("[UserRepo] GetUserByEmailAndPassword called: email='%s'", email)

	if email == "" || passwordHash == "" {
		return nil, errors.New("email and password hash are required")
	}

	query := `SELECT id, email, full_name, role, created_at FROM users WHERE email = $1 AND password_hash = $2`
	row := r.db.QueryRowContext(ctx, query, email, passwordHash)

	var u domain.User
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[UserRepo] GetUserByEmailAndPassword result: invalid credentials for '%s'", email)
			return nil, domain.ErrInvalidCredentials
		}
		log.Printf("[UserRepo] GetUserByEmailAndPassword error: %v", err)
		return nil, fmt.Errorf("database scan failed: %w", err)
	}

	u.CreatedAt = createdAt.Time
	log.Printf("[UserRepo] GetUserByEmailAndPassword success: authenticated user '%s'", u.Email)
	return &u, nil
}

// GetUserByEmail возвращает пользователя и его хеш пароля по email
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	log.Printf("[UserRepo] GetUserByEmail called: email='%s'", email)

	if email == "" {
		return nil, "", errors.New("email cannot be empty")
	}

	query := `SELECT id, email, full_name, role, password_hash, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	var passwordHash string
	var createdAt sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &passwordHash, &createdAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[UserRepo] GetUserByEmail result: user not found")
			return nil, "", domain.ErrUserNotFound
		}
		log.Printf("[UserRepo] GetUserByEmail error: %v", err)
		return nil, "", fmt.Errorf("database scan failed: %w", err)
	}

	u.CreatedAt = createdAt.Time
	log.Printf("[UserRepo] GetUserByEmail success: found user ID=%d", u.ID)
	return &u, passwordHash, nil
}

// UpdatePassword обновляет пароль пользователя
func (r *UserRepo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	log.Printf("[UserRepo] UpdatePassword called: userID=%d", userID)

	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}
	if passwordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, passwordHash, userID)

	if err != nil {
		log.Printf("[UserRepo] UpdatePassword error: %v", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		log.Printf("[UserRepo] UpdatePassword warning: no rows affected for userID=%d", userID)
		return domain.ErrUserNotFound
	}

	log.Printf("[UserRepo] UpdatePassword success: password updated for user ID=%d", userID)
	return nil
}

// UpdateName обновляет только имя пользователя
func (r *UserRepo) UpdateName(ctx context.Context, userID int, name string) error {
	log.Printf("[UserRepo] UpdateName called: userID=%d, name='%s'", userID, name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := r.db.ExecContext(ctx, `UPDATE users SET full_name = $1 WHERE id = $2`, name, userID)
	if err != nil {
		log.Printf("[UserRepo] UpdateName error: %v", err)
		return err
	}
	log.Printf("[UserRepo] UpdateName success")
	return nil
}

func (r *UserRepo) UpdateEmail(ctx context.Context, userID int, email string) error {
	log.Printf("[UserRepo] UpdateEmail called: userID=%d, email='%s'", userID, email)
	if email == "" {
		return errors.New("email cannot be empty")
	}

	query := `UPDATE users SET email = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, email, userID)

	if err != nil {
		// Проверяем код ошибки PostgreSQL: 23505 - нарушение уникального ограничения
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("[UserRepo] UpdateEmail error: duplicate email '%s'", email)
			// Возвращаем специальную ошибку или оборачиваем её
			return fmt.Errorf("email '%s' already exists", email)
		}
		log.Printf("[UserRepo] UpdateEmail error: %v", err)
		return err
	}

	log.Printf("[UserRepo] UpdateEmail success")
	return nil
}

// GetPasswordHash возвращает только хеш пароля
func (r *UserRepo) GetPasswordHash(ctx context.Context, userID int) (string, error) {
	log.Printf("[UserRepo] GetPasswordHash called: userID=%d", userID)
	var hash string
	err := r.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[UserRepo] GetPasswordHash: user not found")
			return "", domain.ErrUserNotFound
		}
		log.Printf("[UserRepo] GetPasswordHash error: %v", err)
		return "", err
	}
	log.Printf("[UserRepo] GetPasswordHash success")
	return hash, nil
}

// SearchByEmail ищет пользователей по подстроке email
func (r *UserRepo) SearchByEmail(ctx context.Context, query string, limit int) ([]*domain.User, error) {
	log.Printf("[UserRepo] SearchByEmail called: query='%s', limit=%d", query, limit)

	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	if limit <= 0 {
		limit = 10 // Значение по умолчанию
	}

	sqlQuery := `SELECT id, email, full_name, role, created_at FROM users WHERE email ILIKE $1 LIMIT $2`
	searchPattern := "%" + query + "%"

	log.Printf("[UserRepo] Executing SQL: %s with args: [%s, %d]", sqlQuery, searchPattern, limit)

	rows, err := r.db.QueryContext(ctx, sqlQuery, searchPattern, limit)
	if err != nil {
		log.Printf("[UserRepo] SearchByEmail query execution failed: %v", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	count := 0

	for rows.Next() {
		u := &domain.User{}
		var fullName string
		var createdAt sql.NullTime

		err := rows.Scan(&u.ID, &u.Email, &fullName, &u.Role, &createdAt)
		if err != nil {
			log.Printf("[UserRepo] SearchByEmail scan failed: %v", err)
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		u.Name = fullName
		if createdAt.Valid {
			u.CreatedAt = createdAt.Time
		}

		users = append(users, u)
		count++
	}

	if err = rows.Err(); err != nil {
		log.Printf("[UserRepo] SearchByEmail rows iteration error: %v", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	log.Printf("[UserRepo] SearchByEmail completed successfully. Found %d users.", count)
	return users, nil
}

// Вспомогательная функция для проверки подстроки (если нужно для обработки ошибок PG)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
