package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type RefreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

// hashToken создает SHA-256 хеш токена
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Save сохраняет хеш refresh-токена в БД
func (r *RefreshTokenRepo) Save(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, hashToken(token), expiresAt)
	return err
}

// Validate проверяет, существует ли токен в БД и не истек ли он
func (r *RefreshTokenRepo) Validate(ctx context.Context, userID int, token string) error {
	query := `
		SELECT id FROM refresh_tokens 
		WHERE user_id = $1 AND token_hash = $2 AND expires_at > NOW()
		LIMIT 1
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, userID, hashToken(token)).Scan(&id)
	if err == sql.ErrNoRows {
		return errors.New("token not found or expired")
	}
	return err
}

// DeleteByUser удаляет все токены пользователя (для Logout)
func (r *RefreshTokenRepo) DeleteByUser(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}
