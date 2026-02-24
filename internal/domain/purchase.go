package domain

import "time"

type Purchase struct {
	UserID      int
	CourseID    int
	PurchasedAt time.Time
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
