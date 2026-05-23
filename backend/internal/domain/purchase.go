package domain

import "time"

type Purchase struct {
	UserID      int       `json:"user_id"`
	CourseID    int       `json:"course_id"`
	PurchasedAt time.Time `json:"purchased_at"`
}

type CourseWithDate struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         int       `json:"price"`
	CoverImageURL string    `json:"cover_image_url"`
	PurchasedAt   time.Time `json:"purchased_at"`
}

func NewPurchase(userID, courseID int) *Purchase {
	return &Purchase{
		UserID:      userID,
		CourseID:    courseID,
		PurchasedAt: time.Now().UTC(),
	}
}

func RestorePurchase(userID, courseID int, purchasedAt time.Time) *Purchase {
	return &Purchase{
		UserID:      userID,
		CourseID:    courseID,
		PurchasedAt: purchasedAt,
	}
}
