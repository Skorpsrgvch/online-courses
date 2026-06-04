package httpintegration

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authhttp "github.com/Skorpsrgvch/online-courses/internal/adapter/http/auth"
	"github.com/Skorpsrgvch/online-courses/internal/domain"
	registeruc "github.com/Skorpsrgvch/online-courses/internal/usecase/auth/register"
	"github.com/gin-gonic/gin"
)

type integrationUserCreator struct {
	user *domain.User
}

func (f *integrationUserCreator) CreateUser(_ context.Context, user *domain.User, _ string) error {
	user.ID = 1
	f.user = user
	return nil
}

func TestRegisterHandlerCreatesUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userCreator := &integrationUserCreator{}
	uc, err := registeruc.NewUsecase(userCreator)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	router := gin.New()
	router.POST("/api/auth/register", authhttp.NewRegisterHandler(uc).Handle)

	body := bytes.NewBufferString(`{"email":"User@Example.COM","full_name":"Test User","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if userCreator.user == nil {
		t.Fatal("expected user to be created")
	}
	if userCreator.user.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", userCreator.user.Email)
	}
}

func TestRegisterHandlerRejectsInvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userCreator := &integrationUserCreator{}
	uc, err := registeruc.NewUsecase(userCreator)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	router := gin.New()
	router.POST("/api/auth/register", authhttp.NewRegisterHandler(uc).Handle)

	body := bytes.NewBufferString(`{"email":"not-an-email","full_name":"Test User","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code == http.StatusCreated {
		t.Fatal("invalid email should not return 201")
	}
	if userCreator.user != nil {
		t.Fatal("user should not be created for invalid email")
	}
}
