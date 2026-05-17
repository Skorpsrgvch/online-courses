package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID        int
	Email     string
	Name      string
	Role      string // "admin" / "user"
	CreatedAt time.Time
}

// NewUser создаёт нового пользователя.
// Валидирует только формат данных, НЕ проверяет существование в БД.
func NewUser(email, name, role string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if role != "admin" && role != "user" {
		return nil, fmt.Errorf("invalid role: %s", role)
	}
	return &User{
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// RestoreUser восстанавливает пользователя из данных БД
func RestoreUser(id int, email, name, role string, createdAt time.Time) *User {
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: createdAt,
	}
}
