package postgres

import (
	"context"
	"database/sql"
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

// CreatePurchase добавляет запись о покупке (например, после фиктивной оплаты)
func (r *PurchaseRepo) CreatePurchase(ctx context.Context, userID, courseID int) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO user_purchases (user_id, course_id) VALUES ($1, $2)",
		userID, courseID,
	)
	return err
}

// GetUserCourseIDs возвращает список ID курсов, купленных пользователем
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
