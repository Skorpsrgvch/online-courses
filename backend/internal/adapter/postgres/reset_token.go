package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"
	"go.uber.org/zap"
)

type ResetTokenRepo struct {
	db *sql.DB
}

func hashResetCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func NewResetTokenRepo(db *sql.DB) *ResetTokenRepo {
	if db == nil {
		zap.L().Fatal("Database connection is nil for ResetTokenRepo")
	}
	return &ResetTokenRepo{db: db}
}

// CreateCode сохраняет код восстановления в БД
func (r *ResetTokenRepo) CreateCode(ctx context.Context, userID int, code string, expiresAt time.Time) error {
	zap.L().Debug("CreateCode called", zap.Int("user_id", userID), zap.Time("expires_at", expiresAt))

	// Инвалидация старых кодов
	queryInvalidate := `UPDATE password_reset_tokens SET used = TRUE WHERE user_id = $1 AND used = FALSE AND expires_at > NOW()`
	res, err := r.db.ExecContext(ctx, queryInvalidate, userID)
	if err != nil {
		zap.L().Error("CreateCode: failed to invalidate old codes", zap.Error(err))
		return fmt.Errorf("failed to invalidate old codes: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		zap.L().Debug("Invalidated old codes", zap.Int("count", int(rowsAffected)), zap.Int("user_id", userID))
	}

	// Вставка нового кода
	queryInsert := `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at, used) VALUES ($1, $2, $3, FALSE)`
	_, err = r.db.ExecContext(ctx, queryInsert, userID, hashResetCode(code), expiresAt)
	if err != nil {
		zap.L().Error("CreateCode: failed to insert new code", zap.Error(err))
		return fmt.Errorf("failed to insert new code: %w", err)
	}

	zap.L().Info("Reset code created", zap.Int("user_id", userID))
	return nil
}

// GetByCode ищет запись по коду
func (r *ResetTokenRepo) GetByCode(ctx context.Context, code string) (*resetpassword.TokenData, error) {
	zap.L().Debug("GetByCode called", zap.String("code", code))

	var data resetpassword.TokenData
	query := `SELECT user_id, expires_at, used FROM password_reset_tokens WHERE token_hash = $1`

	err := r.db.QueryRowContext(ctx, query, hashResetCode(code)).Scan(&data.UserID, &data.ExpiresAt, &data.Used)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Debug("GetByCode: code not found")
			return nil, errors.New("invalid code")
		}
		zap.L().Error("GetByCode: database scan failed", zap.Error(err))
		return nil, fmt.Errorf("database scan failed: %w", err)
	}

	zap.L().Debug("GetByCode success", zap.Int("user_id", data.UserID), zap.Bool("used", data.Used))
	return &data, nil
}

// MarkTokenUsed помечает код как использованный
func (r *ResetTokenRepo) MarkTokenUsed(ctx context.Context, code string) error {
	zap.L().Debug("MarkTokenUsed called", zap.String("code", code))

	query := `UPDATE password_reset_tokens SET used = TRUE WHERE token_hash = $1`
	res, err := r.db.ExecContext(ctx, query, hashResetCode(code))
	if err != nil {
		zap.L().Error("MarkTokenUsed failed", zap.Error(err))
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		zap.L().Warn("MarkTokenUsed: failed to get rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		zap.L().Warn("MarkTokenUsed: no rows updated", zap.String("code", code))
		return errors.New("code not found or already used")
	}

	zap.L().Info("Reset code marked as used", zap.String("code", code))
	return nil
}

// InvalidateOldCodes помечает все коды пользователя как использованные
func (r *ResetTokenRepo) InvalidateOldCodes(ctx context.Context, userID int) error {
	zap.L().Debug("InvalidateOldCodes called", zap.Int("user_id", userID))

	query := `UPDATE password_reset_tokens SET used = TRUE WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		zap.L().Error("InvalidateOldCodes failed", zap.Error(err))
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	zap.L().Info("InvalidateOldCodes completed", zap.Int("rows_affected", int(rowsAffected)))
	return nil
}

// GetLastCodeTime возвращает время создания последнего активного кода
func (r *ResetTokenRepo) GetLastCodeTime(ctx context.Context, userID int) (time.Time, error) {
	var createdAt time.Time

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
