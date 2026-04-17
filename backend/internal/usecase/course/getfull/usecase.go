package getfull

import (
	"context"
	"errors"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type Input struct {
	CourseID int
	UserID   int    // 0 — если неавторизован
	Role     string // "admin", "user"
}

type LessonOutput struct {
	ID             int    `json:"id"`
	ModuleID       int    `json:"module_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	LessonType     string `json:"lesson_type"`
	VideoEmbedID   string `json:"video_embed_id"`
	ArticleContent string `json:"article_content"`
	Order          int    `json:"order"`
}

type ModuleOutput struct {
	ID       int            `json:"id"`
	CourseID int            `json:"course_id"`
	Title    string         `json:"title"`
	Order    int            `json:"order"`
	Lessons  []LessonOutput `json:"lessons"`
}

type Output struct {
	Course  *domain.Course `json:"course"`
	Modules []ModuleOutput `json:"modules"`
}

type CourseReader interface {
	GetByID(ctx context.Context, id int) (*domain.Course, error)
}

type ModuleReader interface {
	GetByCourseID(ctx context.Context, courseID int) ([]*domain.Module, error)
}

type LessonReader interface {
	GetByModuleID(ctx context.Context, moduleID int) ([]*domain.Lesson, error)
}

type PurchaseChecker interface {
	HasPurchased(ctx context.Context, userID, courseID int) (bool, error)
}

type Usecase struct {
	courseReader    CourseReader
	moduleReader    ModuleReader
	lessonReader    LessonReader
	purchaseChecker PurchaseChecker
}

func NewUsecase(courseReader CourseReader, moduleReader ModuleReader, lessonReader LessonReader, purchaseChecker PurchaseChecker) (*Usecase, error) {
	if courseReader == nil || moduleReader == nil || lessonReader == nil || purchaseChecker == nil {
		return nil, errors.New("all dependencies required")
	}
	return &Usecase{
		courseReader:    courseReader,
		moduleReader:    moduleReader,
		lessonReader:    lessonReader,
		purchaseChecker: purchaseChecker,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	course, err := u.courseReader.GetByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	// Бесплатный курс — доступен всем
	if !course.IsPublic {
		if input.UserID == 0 {
			return nil, domain.ErrAccessDenied
		}
		if input.Role != "admin" {
			purchased, err := u.purchaseChecker.HasPurchased(ctx, input.UserID, course.ID)
			if err != nil {
				return nil, err
			}
			if !purchased {
				return nil, domain.ErrCourseNotPurchased
			}
		}
	}

	// Получаем модули
	modules, err := u.moduleReader.GetByCourseID(ctx, course.ID)
	if err != nil {
		return nil, err
	}

	var modulesOut []ModuleOutput
	for _, m := range modules {
		lessons, err := u.lessonReader.GetByModuleID(ctx, m.ID)
		if err != nil {
			return nil, err
		}

		var lessonsOut []LessonOutput
		for _, l := range lessons {
			lessonOut := LessonOutput{
				ID:           l.ID,
				ModuleID:     l.ModuleID,
				Title:        l.Title,
				Description:  l.Description,
				LessonType:   string(l.LessonType),
				Order:        l.Order,
			}
			// Для неоплаченных курсов не возвращаем контент
			if course.IsPublic || input.Role == "admin" {
				lessonOut.VideoEmbedID = l.VideoEmbedID
				lessonOut.ArticleContent = l.ArticleContent
			}
			lessonsOut = append(lessonsOut, lessonOut)
		}

		modulesOut = append(modulesOut, ModuleOutput{
			ID:       m.ID,
			CourseID: m.CourseID,
			Title:    m.Title,
			Order:    m.Order,
			Lessons:  lessonsOut,
		})
	}

	return &Output{
		Course:  course,
		Modules: modulesOut,
	}, nil
}
