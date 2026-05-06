package postgres

import (
	"context"
	"database/sql"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type PurchaseRepo struct {
	db *sql.DB
}

func NewPurchaseRepo(db *sql.DB) *PurchaseRepo {
	return &PurchaseRepo{db: db}
}

// HasPurchased проверяет, купил ли пользователь курс
func (r *PurchaseRepo) HasPurchased(ctx context.Context, userID, courseID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_purchases WHERE user_id = $1 AND course_id = $2)",
		userID, courseID,
	).Scan(&exists)
	return exists, err
}

// Create добавляет запись о покупке
func (r *PurchaseRepo) Create(ctx context.Context, userID, courseID int, paymentID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO user_purchases (user_id, course_id) VALUES ($1, $2)",
		userID, courseID,
	)
	return err
}

// GetUserCourseIDs возвращает список ID курсов
func (r *PurchaseRepo) GetUserCourseIDs(ctx context.Context, userID int) ([]int, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT course_id FROM user_purchases WHERE user_id = $1 ORDER BY purchased_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courseIDs []int
	for rows.Next() {
		var cid int
		if err := rows.Scan(&cid); err != nil {
			return nil, err
		}
		courseIDs = append(courseIDs, cid)
	}
	return courseIDs, nil
}

// GetByUser возвращает полный список покупок с данными курсов
func (r *PurchaseRepo) GetByUser(ctx context.Context, userID int) ([]domain.CourseWithDate, error) {
	// Явно указываем поля courses, которые нам нужны.
	// Если description нет в таблице courses, уберите его из запроса.
	query := `
		SELECT c.id, c.title, COALESCE(c.description, ''), c.price, COALESCE(c.cover_image_url, ''), up.purchased_at
		FROM user_purchases up
		JOIN courses c ON c.id = up.course_id
		WHERE up.user_id = $1
		ORDER BY up.purchased_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []domain.CourseWithDate

	for rows.Next() {
		var p domain.CourseWithDate
		// Порядок полей должен строго соответствовать порядку в SELECT
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.CoverImageURL, &p.PurchasedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}

	// Возвращаем пустой слайс вместо nil, если ничего не найдено (лучше для JSON)
	if purchases == nil {
		return []domain.CourseWithDate{}, nil
	}

	return purchases, nil
}

func (r *PurchaseRepo) GrantAccess(ctx context.Context, userID, courseID int) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO user_purchases (user_id, course_id) 
        VALUES ($1, $2)
        ON CONFLICT (user_id, course_id) DO NOTHING
    `, userID, courseID)
	return err
}
