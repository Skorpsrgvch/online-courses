package getfull

import (
	"context"
	"errors"
)

func NewUsecase(courseReader CourseReader, moduleReader ModuleReader, lessonReader LessonReader, purchaseChecker PurchaseChecker, pr ProgressReader) (*Usecase, error) {
	if courseReader == nil || moduleReader == nil || lessonReader == nil || purchaseChecker == nil {
		return nil, errors.New("all dependencies required")
	}
	return &Usecase{
		courseReader:    courseReader,
		moduleReader:    moduleReader,
		lessonReader:    lessonReader,
		purchaseChecker: purchaseChecker,
		progressReader:  pr,
	}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	course, err := u.courseReader.GetByID(ctx, input.CourseID)
	if err != nil {
		return nil, err
	}

	hasFullAccess := false

	if course.IsPublic {
		hasFullAccess = true
	} else if input.Role == "admin" {
		hasFullAccess = true
	} else if input.UserID > 0 {
		purchased, err := u.purchaseChecker.HasPurchased(ctx, input.UserID, course.ID)
		if err != nil {
			return nil, err
		}
		if purchased {
			hasFullAccess = true
		}
	}

	dbModules, err := u.moduleReader.GetByCourseID(ctx, course.ID)
	if err != nil {
		return nil, err
	}

	modulesOut := make([]ModuleOutput, 0, len(dbModules))

	for _, m := range dbModules {
		dbLessons, err := u.lessonReader.GetByModuleID(ctx, m.ID)
		if err != nil {
			return nil, err
		}

		lessonsOut := make([]LessonOutput, 0, len(dbLessons))

		for _, l := range dbLessons {
			lessonOut := LessonOutput{
				ID:          l.ID,
				ModuleID:    l.ModuleID,
				Title:       l.Title,
				Description: l.Description,
				Order:       l.Order,
			}

			if hasFullAccess {
				lessonOut.VideoEmbedID = l.VideoEmbedID
				lessonOut.PrivateKey = l.PrivateKey

				if input.UserID > 0 {
					isCompleted, _ := u.progressReader.IsLessonCompleted(ctx, input.UserID, l.ID)
					lessonOut.IsCompleted = isCompleted
				}
			} else {
				lessonOut.VideoEmbedID = ""
				lessonOut.PrivateKey = nil
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

	progressPercent := 0
	if input.UserID > 0 && (hasFullAccess || course.IsPublic) {
		completed, total, err := u.progressReader.GetCourseProgress(ctx, input.UserID, course.ID)
		if err == nil && total > 0 {
			progressPercent = int(float64(completed) / float64(total) * 100)
		}
	}

	return &Output{
		Course:          course,
		Modules:         modulesOut,
		IsPurchased:     hasFullAccess,
		ProgressPercent: progressPercent,
	}, nil
}
