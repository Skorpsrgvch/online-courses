package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"go.uber.org/zap"
)

type ReviewRepo struct {
	db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) CreateReview(ctx context.Context, review *domain.Review) error {
	query := `
		INSERT INTO reviews (user_id, course_id, text, rating, approved, rejection_reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET 
			text = EXCLUDED.text,
			rating = EXCLUDED.rating,
			approved = false,
			rejection_reason = NULL,
			created_at = COALESCE(EXCLUDED.created_at, reviews.created_at)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		review.UserID,
		review.CourseID,
		review.Text,
		review.Rating,
		review.Approved,
		review.RejectionReason,
		review.CreatedAt,
	).Scan(&review.ID)

	if err != nil {
		zap.L().Error("ReviewRepo.CreateReview failed",
			zap.Int("user_id", review.UserID),
			zap.Int("course_id", review.CourseID),
			zap.Error(err),
		)
		return err
	}

	zap.L().Debug("Review created/updated", zap.Int("id", review.ID))
	return nil
}

func (r *ReviewRepo) GetApprovedReviewsByCourse(ctx context.Context, courseID int) ([]*domain.Review, error) {
	query := `
        SELECT 
            r.id, r.user_id, r.course_id, r.text, r.rating, r.approved, 
            COALESCE(r.rejection_reason, ''), r.created_at, 
            COALESCE(NULLIF(TRIM(u.full_name), ''), 'Аноним')
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        WHERE r.course_id = $1 AND r.approved = true
        ORDER BY r.created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var rev domain.Review
		var authorName string
		var rejectionReason string

		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.CourseID, &rev.Text, &rev.Rating,
			&rev.Approved, &rejectionReason, &rev.CreatedAt, &authorName,
		)
		if err != nil {
			return nil, err
		}
		rev.RejectionReason = rejectionReason
		rev.AuthorName = authorName
		reviews = append(reviews, &rev)
	}

	return reviews, rows.Err()
}

func (r *ReviewRepo) GetByUserAndCourse(ctx context.Context, userID, courseID int) (*domain.Review, error) {
	query := `
		SELECT 
			r.id, r.user_id, r.course_id, r.text, r.rating, r.approved,
			COALESCE(r.rejection_reason, ''), r.created_at,
			COALESCE(NULLIF(TRIM(u.full_name), ''), 'Аноним')
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.user_id = $1 AND r.course_id = $2
		LIMIT 1
	`

	var rev domain.Review
	var authorName, rejectionReason string

	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&rev.ID, &rev.UserID, &rev.CourseID, &rev.Text, &rev.Rating,
		&rev.Approved, &rejectionReason, &rev.CreatedAt, &authorName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rev.RejectionReason = rejectionReason
	rev.AuthorName = authorName
	return &rev, nil
}

func (r *ReviewRepo) GetPendingReviews(ctx context.Context) ([]*domain.Review, error) {
	query := `
        SELECT 
            r.id, r.user_id, r.course_id, r.rating, r.text, r.approved,
            COALESCE(r.rejection_reason, ''), r.created_at,
            u.full_name as author_name, c.title as course_title
        FROM reviews r
        JOIN users u ON r.user_id = u.id
        JOIN courses c ON r.course_id = c.id
        WHERE (r.approved IS NOT TRUE) AND (r.rejection_reason IS NULL OR r.rejection_reason = '')
        ORDER BY r.created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var rev domain.Review
		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.CourseID, &rev.Rating, &rev.Text,
			&rev.Approved, &rev.RejectionReason, &rev.CreatedAt, &rev.AuthorName, &rev.CourseTitle,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &rev)
	}

	return reviews, rows.Err()
}

func (r *ReviewRepo) GetByUserID(ctx context.Context, userID int) ([]*domain.Review, error) {
	query := `
		SELECT 
			r.id, r.user_id, r.course_id, r.rating, r.text, r.approved,
			COALESCE(r.rejection_reason, ''), r.created_at,
			COALESCE(NULLIF(TRIM(u.full_name), ''), 'Аноним'), c.title
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		JOIN courses c ON r.course_id = c.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var rev domain.Review
		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.CourseID, &rev.Rating, &rev.Text,
			&rev.Approved, &rev.RejectionReason, &rev.CreatedAt, &rev.AuthorName, &rev.CourseTitle,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &rev)
	}

	return reviews, rows.Err()
}

func (r *ReviewRepo) ApproveReview(ctx context.Context, reviewID int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE reviews SET approved = true, rejection_reason = NULL WHERE id = $1`, reviewID)
	return err
}

func (r *ReviewRepo) RejectReview(ctx context.Context, reviewID int, reason string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE reviews SET approved = false, rejection_reason = $2 WHERE id = $1`, reviewID, reason)
	return err
}

func (r *ReviewRepo) DeleteReview(ctx context.Context, reviewID int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, reviewID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("review not found or access denied")
	}
	return nil
}

func (r *ReviewRepo) DeleteByAdmin(ctx context.Context, reviewID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, reviewID)
	return err
}
