package domain

import (
	"errors"
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusSucceeded         PaymentStatus = "succeeded"
	PaymentStatusCanceled          PaymentStatus = "canceled"
	PaymentStatusFailed            PaymentStatus = "failed"
	PaymentStatusWaitingForCapture PaymentStatus = "waiting_for_capture"
)

// Payment represents a payment entity
type Payment struct {
	ID              int               `json:"id"`
	PaymentID       string            `json:"payment_id"` // UUID from YooKassa
	UserID          int               `json:"user_id"`
	CourseID        int               `json:"course_id"`
	Amount          int               `json:"amount"` // in kopecks
	Currency        string            `json:"currency"`
	Status          PaymentStatus     `json:"status"`
	ConfirmationURL string            `json:"confirmation_url,omitempty"`
	Description     string            `json:"description,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	ExpiresAt       *time.Time        `json:"expires_at,omitempty"`
	PaidAt          *time.Time        `json:"paid_at,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	UserID      int    `json:"user_id"`
	CourseID    int    `json:"course_id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	ReturnURL   string `json:"return_url"`
}

// PaymentConfirmation represents confirmation data for payment
type PaymentConfirmation struct {
	Type              string `json:"type"`
	ConfirmationToken string `json:"confirmation_token,omitempty"`
	Enforce           bool   `json:"enforce,omitempty"`
	ReturnURL         string `json:"return_url,omitempty"`
	ConfirmationURL   string `json:"confirmation_url,omitempty"`
	ConfirmationData  string `json:"confirmation_data,omitempty"`
	Locale            string `json:"locale,omitempty"`
}

// PaymentErrors
var (
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentExpired        = errors.New("payment expired")
	ErrPaymentAlreadyPaid    = errors.New("payment already paid")
	ErrPaymentInvalidStatus  = errors.New("invalid payment status")
	ErrPaymentCreationFailed = errors.New("failed to create payment")
)

// NewPayment creates a new payment entity
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

// IsExpired checks if the payment has expired
func (p *Payment) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*p.ExpiresAt)
}

// CanBeCompleted checks if the payment can be completed based on its status
func (p *Payment) CanBeCompleted() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusWaitingForCapture
}

// SetSucceeded marks the payment as succeeded
func (p *Payment) SetSucceeded() {
	now := time.Now().UTC()
	p.Status = PaymentStatusSucceeded
	p.PaidAt = &now
}

// SetCanceled marks the payment as canceled
func (p *Payment) SetCanceled() {
	p.Status = PaymentStatusCanceled
}

// SetFailed marks the payment as failed
func (p *Payment) SetFailed() {
	p.Status = PaymentStatusFailed
}
