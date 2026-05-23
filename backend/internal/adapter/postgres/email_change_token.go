package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type EmailChangeTokenData struct {
	UserID    int
	NewEmail  string
	ExpiresAt time.Time
	Used      bool
}

type EmailChangeTokenRepo struct {
	db *sql.DB
}

func NewEmailChangeTokenRepo(db *sql.DB) *EmailChangeTokenRepo {
	return &EmailChangeTokenRepo{db: db}
}

func (r *EmailChangeTokenRepo) CreateToken(ctx context.Context, userID int, newEmail, token string, expiresAt time.Time) error {
	query := `INSERT INTO email_change_tokens (user_id, new_email, token, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, userID, newEmail, token, expiresAt)
	if err != nil {
		zap.L().Error("Failed to create email change token", zap.Int("userID", userID), zap.Error(err))
		return fmt.Errorf("create token: %w", err)
	}
	zap.L().Debug("Email change token created", zap.Int("userID", userID), zap.Time("expiresAt", expiresAt))
	return nil
}

func (r *EmailChangeTokenRepo) GetToken(ctx context.Context, token string) (*EmailChangeTokenData, error) {
	var data EmailChangeTokenData
	query := `SELECT user_id, new_email, expires_at, used FROM email_change_tokens WHERE token = $1`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&data.UserID, &data.NewEmail, &data.ExpiresAt, &data.Used,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email change token")
		}
		zap.L().Error("Failed to get email change token", zap.Error(err))
		return nil, fmt.Errorf("get token: %w", err)
	}
	return &data, nil
}

func (r *EmailChangeTokenRepo) MarkTokenUsed(ctx context.Context, token string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE email_change_tokens SET used = TRUE WHERE token = $1`, token)
	if err != nil {
		zap.L().Error("Failed to mark token as used", zap.Error(err))
		return fmt.Errorf("mark token used: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("token not found or already used")
	}

	zap.L().Debug("Email change token marked as used", zap.String("token", token))
	return nil
}

func (r *EmailChangeTokenRepo) DeleteExpired(ctx context.Context) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM email_change_tokens WHERE expires_at < NOW()`)
	if err != nil {
		zap.L().Error("Failed to delete expired tokens", zap.Error(err))
		return fmt.Errorf("delete expired tokens: %w", err)
	}

	count, _ := res.RowsAffected()
	if count > 0 {
		zap.L().Info("Deleted expired email change tokens", zap.Int64("count", count))
	}
	return nil
}
