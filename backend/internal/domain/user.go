package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(email, name, role string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required: %w", ErrInvalidInput)
	}
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", ErrInvalidInput)
	}
	if role != "admin" && role != "user" {
		return nil, fmt.Errorf("invalid role '%s': %w", role, ErrInvalidInput)
	}
	return &User{
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func RestoreUser(id int, email, name, role string, createdAt time.Time) *User {
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: createdAt,
	}
}
