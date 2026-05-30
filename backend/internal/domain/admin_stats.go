package domain

import "time"

type StudentStat struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	RegisteredAt     time.Time  `json:"registered_at"`
	TotalCourses     int        `json:"total_courses"`
	CompletedCourses int        `json:"completed_courses"`
	TotalLessons     int        `json:"total_lessons"`
	CompletedLessons int        `json:"completed_lessons"`
	ProgressPercent  int        `json:"progress_percent"`
	LastActivityAt   *time.Time `json:"last_activity_at"`
}

type StudentCourseDetail struct {
	CourseID         int        `json:"course_id"`
	CourseTitle      string     `json:"course_title"`
	PurchasedAt      time.Time  `json:"purchased_at"`
	TotalLessons     int        `json:"total_lessons"`
	CompletedLessons int        `json:"completed_lessons"`
	ProgressPercent  int        `json:"progress_percent"`
	LastLessonTitle  string     `json:"last_lesson_title"`
	LastActivityAt   *time.Time `json:"last_activity_at"`
}
