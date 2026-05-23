package domain

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusSucceeded         PaymentStatus = "succeeded"
	PaymentStatusCanceled          PaymentStatus = "canceled"
	PaymentStatusFailed            PaymentStatus = "failed"
	PaymentStatusWaitingForCapture PaymentStatus = "waiting_for_capture"
)

type Payment struct {
	ID              int               `json:"id"`
	PaymentID       string            `json:"payment_id"`
	UserID          int               `json:"user_id"`
	CourseID        int               `json:"course_id"`
	Amount          int               `json:"amount"`
	Currency        string            `json:"currency"`
	Status          PaymentStatus     `json:"status"`
	ConfirmationURL string            `json:"confirmation_url,omitempty"`
	Description     string            `json:"description,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	ExpiresAt       *time.Time        `json:"expires_at,omitempty"`
	PaidAt          *time.Time        `json:"paid_at,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type CreatePaymentRequest struct {
	UserID      int    `json:"user_id"`
	CourseID    int    `json:"course_id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	ReturnURL   string `json:"return_url"`
}

type PaymentConfirmation struct {
	Type              string `json:"type"`
	ConfirmationToken string `json:"confirmation_token,omitempty"`
	Enforce           bool   `json:"enforce,omitempty"`
	ReturnURL         string `json:"return_url,omitempty"`
	ConfirmationURL   string `json:"confirmation_url,omitempty"`
	ConfirmationData  string `json:"confirmation_data,omitempty"`
	Locale            string `json:"locale,omitempty"`
}

func NewPayment(paymentID string, userID, courseID, amount int, currency, description string) *Payment {
	now := time.Now().UTC()
	return &Payment{
		PaymentID:   paymentID,
		UserID:      userID,
		CourseID:    courseID,
		Amount:      amount,
		Currency:    currency,
		Status:      PaymentStatusPending,
		CreatedAt:   now,
		Description: description,
		Metadata:    make(map[string]string),
	}
}

func (p *Payment) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*p.ExpiresAt)
}

func (p *Payment) CanBeCompleted() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusWaitingForCapture
}

func (p *Payment) SetSucceeded() {
	now := time.Now().UTC()
	p.Status = PaymentStatusSucceeded
	p.PaidAt = &now
}

func (p *Payment) SetCanceled() {
	p.Status = PaymentStatusCanceled
}

func (p *Payment) SetFailed() {
	p.Status = PaymentStatusFailed
}
