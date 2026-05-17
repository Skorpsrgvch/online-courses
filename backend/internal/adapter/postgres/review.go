package postgres

import (
	"context"
	"database/sql"
	"log"
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
	log.Printf("[DEBUG] ReviewRepo.CreateReview: UserID=%d, CourseID=%d, Rating=%d, Text=%q",
		review.UserID, review.CourseID, review.Rating, review.Text)

	query := `
		INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET 
			text = EXCLUDED.text,
			rating = EXCLUDED.rating,
			approved = EXCLUDED.approved,
			created_at = COALESCE(EXCLUDED.created_at, reviews.created_at)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		review.UserID,
		review.CourseID,
		review.Text,
		review.Rating,
		review.Approved,
		review.CreatedAt,
	).Scan(&review.ID)

	if err != nil {
		log.Printf("[ERROR] ReviewRepo.CreateReview: DB error: %v", err)
		return err
	}

	log.Printf("[INFO] ReviewRepo.CreateReview: Success, ID=%d", review.ID)
	return nil
}

func (r *ReviewRepo) GetApprovedReviewsByCourse(ctx context.Context, courseID int) ([]*domain.Review, error) {

	query := `
        SELECT 
    r.id, 
    r.user_id, 
    r.course_id, 
    r.text, 
    r.rating, 
    r.approved, 
    r.created_at, 
    COALESCE(NULLIF(TRIM(u.full_name), ''), 'Аноним') as author_name
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

	reviews := make([]*domain.Review, 0)

	for rows.Next() {
		var (
			id, userID, courseIDDB int
			text                   string
			rating                 int
			approved               bool
			createdAt              time.Time
			authorName             string
		)
		// Добавляем authorName в Scan
		err := rows.Scan(&id, &userID, &courseIDDB, &text, &rating, &approved, &createdAt, &authorName)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Row scanned: ID=%d, UserID=%d, AuthorName='%s'", id, userID, authorName)

		reviews = append(reviews, domain.RestoreReview(id, userID, courseIDDB, text, rating, approved, createdAt, authorName))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepo) GetPendingReviews(ctx context.Context) ([]*domain.Review, error) {
	// ИСПРАВЛЕНО: выбираем approved вместо status, так как в БД и структуре используется approved
	query := `
        SELECT 
            r.id, r.user_id, r.course_id, r.rating, r.text, r.approved, r.created_at,
            u.full_name as author_name,
            c.title as course_title
        FROM reviews r
        JOIN users u ON r.user_id = u.id
        JOIN courses c ON r.course_id = c.id
        WHERE r.approved = false  -- Ищем неодобренные отзывы
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
			&rev.Approved, &rev.CreatedAt, &rev.AuthorName, &rev.CourseTitle,
		)
		if err != nil {
			return nil, err
		}
		log.Printf("[ReviewRepo] GetPending: Found review ID=%d | CourseID=%d | CourseTitle='%s' | Author='%s' | Rating=%d",
			rev.ID, rev.CourseID, rev.CourseTitle, rev.AuthorName, rev.Rating)
		reviews = append(reviews, &rev)
	}

	return reviews, rows.Err()
}

func (r *ReviewRepo) ApproveReview(ctx context.Context, reviewID int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE reviews SET approved = true WHERE id = $1`, reviewID)
	return err
}

func (r *ReviewRepo) RejectReview(ctx context.Context, reviewID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, reviewID)
	return err
}
