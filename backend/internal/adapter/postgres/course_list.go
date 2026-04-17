package postgres

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

// ListAll возвращает список всех активных курсов
func (r *CourseRepo) ListAll(ctx context.Context) ([]*domain.Course, error) {
	query := `
		SELECT id, title, description, is_public, price, author_id, is_active,
		       COALESCE(cover_image_url, '')
		FROM courses
		WHERE is_active = true
		ORDER BY id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.IsPublic,
			&c.Price,
			&c.AuthorID,
			&c.IsActive,
			&c.CoverImageURL,
		); err != nil {
			return nil, err
		}
		courses = append(courses, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
