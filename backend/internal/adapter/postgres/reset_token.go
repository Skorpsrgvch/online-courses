package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/usecase/auth/resetpassword"
)

type ResetTokenRepo struct {
	db *sql.DB
}

func NewResetTokenRepo(db *sql.DB) *ResetTokenRepo {
	return &ResetTokenRepo{db: db}
}

func (r *ResetTokenRepo) CreateResetToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *ResetTokenRepo) GetResetToken(ctx context.Context, token string) (*resetpassword.TokenData, error) {
	var data resetpassword.TokenData
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, expires_at, used FROM password_reset_tokens WHERE token = $1`,
		token,
	).Scan(&data.UserID, &data.ExpiresAt, &data.Used)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid reset token")
		}
		return nil, err
	}
	return &data, nil
}

func (r *ResetTokenRepo) MarkTokenUsed(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE password_reset_tokens SET used = TRUE WHERE token = $1`,
		token,
	)
	return err
}
