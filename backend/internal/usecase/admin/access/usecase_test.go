package access

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type fakeAccessUserRepo struct {
	user *domain.User
	err  error
}

func (f *fakeAccessUserRepo) GetUserByID(_ context.Context, _ int) (*domain.User, error) {
	return f.user, f.err
}

type fakeAccessCourseRepo struct {
	course *domain.Course
	err    error
}

func (f *fakeAccessCourseRepo) GetByID(_ context.Context, _ int) (*domain.Course, error) {
	return f.course, f.err
}

type fakeAccessPurchaseRepo struct {
	called bool
	err    error
}

func (f *fakeAccessPurchaseRepo) GrantAccess(_ context.Context, _, _ int) error {
	f.called = true
	return f.err
}

type fakeAccessEmailSender struct {
	called      bool
	to          string
	userName    string
	courseTitle string
	err         error
}

func (f *fakeAccessEmailSender) SendCourseAccessGranted(to, userName, courseTitle string) error {
	f.called = true
	f.to = to
	f.userName = userName
	f.courseTitle = courseTitle
	return f.err
}

func TestExecuteGrantsAccessAndSendsEmail(t *testing.T) {
	emailSender := &fakeAccessEmailSender{}
	purchaseRepo := &fakeAccessPurchaseRepo{}
	uc, err := NewUsecase(
		&fakeAccessUserRepo{user: domain.RestoreUser(7, "user@example.com", "User", "user", time.Now())},
		&fakeAccessCourseRepo{course: domain.RestoreCourse(3, 1, "Course", "", false, 1000, true, "", time.Now(), "", "", "", "", "", nil)},
		purchaseRepo,
		emailSender,
	)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	output, err := uc.Execute(context.Background(), Input{UserID: 7, CourseID: 3})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if output == nil || !output.Success {
		t.Fatalf("expected successful output, got %#v", output)
	}
	if !purchaseRepo.called {
		t.Fatal("expected purchase repo to be called")
	}
	if !emailSender.called {
		t.Fatal("expected email sender to be called")
	}
	if emailSender.to != "user@example.com" || emailSender.userName != "User" || emailSender.courseTitle != "Course" {
		t.Fatalf("unexpected email args: %#v", emailSender)
	}
}

func TestExecuteDoesNotFailWhenEmailSendingFails(t *testing.T) {
	emailSender := &fakeAccessEmailSender{err: errors.New("smtp failed")}
	uc, err := NewUsecase(
		&fakeAccessUserRepo{user: domain.RestoreUser(7, "user@example.com", "User", "user", time.Now())},
		&fakeAccessCourseRepo{course: domain.RestoreCourse(3, 1, "Course", "", false, 1000, true, "", time.Now(), "", "", "", "", "", nil)},
		&fakeAccessPurchaseRepo{},
		emailSender,
	)
	if err != nil {
		t.Fatalf("NewUsecase returned error: %v", err)
	}

	output, err := uc.Execute(context.Background(), Input{UserID: 7, CourseID: 3})
	if err != nil {
		t.Fatalf("email failure should not fail access grant: %v", err)
	}
	if output == nil || !output.Success {
		t.Fatalf("expected successful output, got %#v", output)
	}
	if !emailSender.called {
		t.Fatal("expected email sender to be called")
	}
}
