package register

import (
	"context"
	"strings"
	"testing"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserCreator struct {
	user         *domain.User
	passwordHash string
	err          error
}

func (f *fakeUserCreator) CreateUser(_ context.Context, user *domain.User, passwordHash string) error {
	if f.err != nil {
		return f.err
	}
	f.user = user
	f.passwordHash = passwordHash
	return nil
}

func TestExecuteCreatesUserWithNormalizedEmailAndHashedPassword(t *testing.T) {
	repo := &fakeUserCreator{}
	uc, err := NewUsecase(repo)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	err = uc.Execute(context.Background(), Input{
		Email:    "  USER@Example.COM ",
		FullName: " Test User ",
		Password: "secret123",
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if repo.user == nil {
		t.Fatal("expected user to be created")
	}
	if repo.user.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", repo.user.Email)
	}
	if repo.user.Name != "Test User" {
		t.Fatalf("expected trimmed full name, got %q", repo.user.Name)
	}
	if repo.passwordHash == "" || repo.passwordHash == "secret123" {
		t.Fatal("expected password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.passwordHash), []byte("secret123")); err != nil {
		t.Fatalf("password hash does not match original password: %v", err)
	}
}

func TestExecuteRejectsInvalidEmail(t *testing.T) {
	repo := &fakeUserCreator{}
	uc, err := NewUsecase(repo)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	err = uc.Execute(context.Background(), Input{
		Email:    "not-an-email",
		FullName: "Test User",
		Password: "secret123",
		Role:     "user",
	})
	if err == nil {
		t.Fatal("expected invalid email error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "email") {
		t.Fatalf("expected email error, got %v", err)
	}
	if repo.user != nil {
		t.Fatal("user should not be created for invalid email")
	}
}

func TestExecuteRejectsWeakPassword(t *testing.T) {
	repo := &fakeUserCreator{}
	uc, err := NewUsecase(repo)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	err = uc.Execute(context.Background(), Input{
		Email:    "user@example.com",
		FullName: "Test User",
		Password: "123",
		Role:     "user",
	})
	if err == nil {
		t.Fatal("expected weak password error")
	}
	if repo.user != nil {
		t.Fatal("user should not be created for weak password")
	}
}
