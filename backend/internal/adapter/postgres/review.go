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

func (r *ReviewRepo) CreateReview(ctx context.Context, review *domain.Review) error {
	query := `
		INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET 
			text = EXCLUDED.text,
			rating = EXCLUDED.rating,
			approved = EXCLUDED.approved, -- или false, если нужна премодерация при изменении
			created_at = COALESCE(EXCLUDED.created_at, reviews.created_at) -- сохраняем старую дату или новую
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		review.UserID,
		review.CourseID,
		review.Text,
		review.Rating,
		review.Approved, // обычно false для новых/измененных
		review.CreatedAt,
	).Scan(&review.ID)
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

	// ИСПРАВЛЕНИЕ: Инициализируем срез, чтобы он не был nil
	reviews := make([]*domain.Review, 0)

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

	// Проверка ошибок итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepo) GetPendingReviews(ctx context.Context) ([]*domain.Review, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, course_id, text, rating, approved, created_at
		 FROM reviews
		 WHERE approved = false
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*domain.Review, 0)

	for rows.Next() {
		var (
			id, userID, courseID int
			text                 string
			rating               int
			approved             bool
			createdAt            time.Time
		)
		err := rows.Scan(&id, &userID, &courseID, &text, &rating, &approved, &createdAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, domain.RestoreReview(id, userID, courseID, text, rating, approved, createdAt))
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

func (r *ReviewRepo) RejectReview(ctx context.Context, reviewID int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM reviews WHERE id = $1`,
		reviewID,
	)
	return err
}
