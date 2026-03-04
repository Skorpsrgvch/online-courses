package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ReviewRepo struct {
	db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) GetApprovedReviewsByCourse(ctx context.Context, courseID int) ([]*domain.Review, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, course_id, text, rating, approved, created_at
		 FROM reviews
		 WHERE course_id = $1 AND approved = true
		 ORDER BY created_at DESC`,
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var (
			id, userID, courseIDDB int
			text                   string
			rating                 int
			approved               bool
			createdAt              time.Time
		)
		err := rows.Scan(&id, &userID, &courseIDDB, &text, &rating, &approved, &createdAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, domain.RestoreReview(id, userID, courseIDDB, text, rating, approved, createdAt))
	}
	return reviews, nil
}

func (r *ReviewRepo) ApproveReview(ctx context.Context, reviewID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE reviews SET approved = true WHERE id = $1`,
		reviewID,
	)
	return err
}
