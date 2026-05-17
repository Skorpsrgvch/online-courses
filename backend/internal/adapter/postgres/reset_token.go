package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"
)

type ResetTokenRepo struct {
	db *sql.DB
}

func NewResetTokenRepo(db *sql.DB) *ResetTokenRepo {
	if db == nil {
		log.Fatal("Database connection is nil for ResetTokenRepo")
	}
	return &ResetTokenRepo{db: db}
}

// CreateCode сохраняет код восстановления в БД
func (r *ResetTokenRepo) CreateCode(ctx context.Context, userID int, code string, expiresAt time.Time) error {
	log.Printf("[ResetTokenRepo] CreateCode called: userID=%d, expiresAt=%s", userID, expiresAt.Format(time.RFC3339))

	// Инвалидация старых кодов
	queryInvalidate := `UPDATE password_reset_tokens SET used = TRUE WHERE user_id = $1 AND used = FALSE AND expires_at > NOW()`
	res, err := r.db.ExecContext(ctx, queryInvalidate, userID)
	if err != nil {
		log.Printf("[ResetTokenRepo] CreateCode error (invalidating old codes): %v", err)
		return fmt.Errorf("failed to invalidate old codes: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	log.Printf("[ResetTokenRepo] Invalidated %d old codes for user %d", rowsAffected, userID)

	// Вставка нового кода
	queryInsert := `INSERT INTO password_reset_tokens (user_id, token, expires_at, used) VALUES ($1, $2, $3, FALSE)`
	_, err = r.db.ExecContext(ctx, queryInsert, userID, code, expiresAt)
	if err != nil {
		log.Printf("[ResetTokenRepo] CreateCode error (inserting new code): %v", err)
		return fmt.Errorf("failed to insert new code: %w", err)
	}

	log.Printf("[ResetTokenRepo] CreateCode success: code created for user %d", userID)
	return nil
}

// GetByCode ищет запись по коду
func (r *ResetTokenRepo) GetByCode(ctx context.Context, code string) (*resetpassword.TokenData, error) {
	log.Printf("[ResetTokenRepo] GetByCode called: code=%s", code)

	var data resetpassword.TokenData
	query := `SELECT user_id, expires_at, used FROM password_reset_tokens WHERE token = $1`

	err := r.db.QueryRowContext(ctx, query, code).Scan(&data.UserID, &data.ExpiresAt, &data.Used)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[ResetTokenRepo] GetByCode result: code not found")
			return nil, errors.New("invalid code")
		}
		log.Printf("[ResetTokenRepo] GetByCode database error: %v", err)
		return nil, fmt.Errorf("database scan failed: %w", err)
	}

	log.Printf("[ResetTokenRepo] GetByCode success: found code for user %d, used=%v, expires=%s",
		data.UserID, data.Used, data.ExpiresAt.Format(time.RFC3339))
	return &data, nil
}

// MarkTokenUsed помечает код как использованный
func (r *ResetTokenRepo) MarkTokenUsed(ctx context.Context, code string) error {
	log.Printf("[ResetTokenRepo] MarkTokenUsed called: code=%s", code)

	query := `UPDATE password_reset_tokens SET used = TRUE WHERE token = $1`
	res, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		log.Printf("[ResetTokenRepo] MarkTokenUsed error: %v", err)
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("[ResetTokenRepo] MarkTokenUsed error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("[ResetTokenRepo] MarkTokenUsed warning: no rows updated for code %s", code)
		return errors.New("code not found or already used")
	}

	log.Printf("[ResetTokenRepo] MarkTokenUsed success: code %s marked as used", code)
	return nil
}

// InvalidateOldCodes (опционально)
func (r *ResetTokenRepo) InvalidateOldCodes(ctx context.Context, userID int) error {
	log.Printf("[ResetTokenRepo] InvalidateOldCodes called: userID=%d", userID)

	query := `UPDATE password_reset_tokens SET used = TRUE WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Printf("[ResetTokenRepo] InvalidateOldCodes error: %v", err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	log.Printf("[ResetTokenRepo] InvalidateOldCodes success: %d rows updated", rowsAffected)
	return nil
}

func (r *ResetTokenRepo) GetLastCodeTime(ctx context.Context, userID int) (time.Time, error) {
	var createdAt time.Time

	// Берем самый свежий неиспользованный токен
	query := `SELECT created_at FROM password_reset_tokens 
	          WHERE user_id = $1 AND used = FALSE AND expires_at > NOW() 
	          ORDER BY created_at DESC LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil // Нет активных кодов
		}
		return time.Time{}, err
	}

	return createdAt, nil
}
