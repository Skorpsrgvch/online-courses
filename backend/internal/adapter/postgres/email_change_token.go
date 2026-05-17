package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

// CreateToken создает токен для смены email
func (r *EmailChangeTokenRepo) CreateToken(ctx context.Context, userID int, newEmail, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO email_change_tokens (user_id, new_email, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, userID, newEmail, token, expiresAt)
	return err
}

// GetToken проверяет токен и возвращает данные
func (r *EmailChangeTokenRepo) GetToken(ctx context.Context, token string) (*EmailChangeTokenData, error) {
	var data EmailChangeTokenData
	query := `
		SELECT user_id, new_email, expires_at, used 
		FROM email_change_tokens 
		WHERE token = $1
	`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&data.UserID, &data.NewEmail, &data.ExpiresAt, &data.Used,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email change token")
		}
		return nil, err
	}
	return &data, nil
}

// MarkTokenUsed помечает токен как использованный
func (r *EmailChangeTokenRepo) MarkTokenUsed(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE email_change_tokens SET used = TRUE WHERE token = $1`,
		token,
	)
	return err
}

// DeleteExpired удаляет старые токены (можно запускать по крону)
func (r *EmailChangeTokenRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM email_change_tokens WHERE expires_at < NOW()`,
	)
	return err
}
