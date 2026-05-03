package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"github.com/google/uuid"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, user_id, course_id, amount, currency, status, confirmation_url, description, created_at, expires_at, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		uuid.MustParse(payment.PaymentID),
		payment.UserID,
		payment.CourseID,
		payment.Amount,
		payment.Currency,
		string(payment.Status),
		payment.ConfirmationURL,
		payment.Description,
		payment.CreatedAt,
		payment.ExpiresAt,
		payment.PaidAt,
	)
	return err
}

func (r *PaymentRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	query := `
		SELECT id, user_id, course_id, amount, currency, status, confirmation_url, description, created_at, expires_at, paid_at
		FROM payments
		WHERE id = $1
	`
	p := &domain.Payment{}
	var statusStr string

	err := r.db.QueryRowContext(ctx, query, uuid.MustParse(paymentID)).Scan(
		&p.PaymentID,
		&p.UserID,
		&p.CourseID,
		&p.Amount,
		&p.Currency,
		&statusStr,
		&p.ConfirmationURL,
		&p.Description,
		&p.CreatedAt,
		&p.ExpiresAt,
		&p.PaidAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	p.Status = domain.PaymentStatus(statusStr)
	return p, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments 
		SET status = $1, paid_at = $2, confirmation_url = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		string(payment.Status),
		payment.PaidAt,
		payment.ConfirmationURL,
		uuid.MustParse(payment.PaymentID),
	)
	return err
}

func (r *PaymentRepository) UpdateConfirmationURL(ctx context.Context, paymentID string, url string) error {
	query := `UPDATE payments SET confirmation_url = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, url, uuid.MustParse(paymentID))
	return err
}

func (r *PaymentRepository) SetExpiresAt(ctx context.Context, paymentID string, expiresAt time.Time) error {
	query := `UPDATE payments SET expires_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, expiresAt, uuid.MustParse(paymentID))
	return err
}

func (r *PaymentRepository) ListByUser(ctx context.Context, userID int) ([]*domain.Payment, error) {
	query := `
		SELECT id, user_id, course_id, amount, currency, status, confirmation_url, description, created_at, expires_at, paid_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		p := &domain.Payment{}
		var statusStr string
		err := rows.Scan(
			&p.PaymentID,
			&p.UserID,
			&p.CourseID,
			&p.Amount,
			&p.Currency,
			&statusStr,
			&p.ConfirmationURL,
			&p.Description,
			&p.CreatedAt,
			&p.ExpiresAt,
			&p.PaidAt,
		)
		if err != nil {
			return nil, err
		}
		p.Status = domain.PaymentStatus(statusStr)
		payments = append(payments, p)
	}
	return payments, nil
}
