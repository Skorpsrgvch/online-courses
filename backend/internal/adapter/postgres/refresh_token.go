package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"go.uber.org/zap"
)

type RefreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

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
	if err != nil {
		zap.L().Error("Failed to save refresh token", zap.Int("user_id", userID), zap.Error(err))
		return err
	}
	zap.L().Debug("Refresh token saved", zap.Int("user_id", userID))
	return nil
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

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("token not found or expired")
		}
		zap.L().Error("Failed to validate refresh token", zap.Int("user_id", userID), zap.Error(err))
		return err
	}
	return nil
}

// DeleteByUser удаляет все токены пользователя (для Logout)
func (r *RefreshTokenRepo) DeleteByUser(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	if err != nil {
		zap.L().Error("Failed to delete refresh tokens", zap.Int("user_id", userID), zap.Error(err))
		return err
	}
	zap.L().Info("All refresh tokens deleted for user", zap.Int("user_id", userID))
	return nil
}
