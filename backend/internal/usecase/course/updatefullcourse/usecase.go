package updatefullcourse

import (
	"context"
)

type Usecase struct {
	repo CourseRepository
}

func NewUsecase(repo CourseRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	return u.repo.UpdateFullWithModules(ctx, input)
}

// ReorderModules принимает список ID модулей в том порядке, в котором они должны идти
func (u *Usecase) ReorderModules(ctx context.Context, courseID int, moduleIDs []int) error {
	orderMap := make(map[int]int)
	for i, id := range moduleIDs {
		orderMap[id] = i + 1 // Порядок обычно с 1
	}
	return u.repo.ReorderModules(ctx, courseID, orderMap)
}

// ReorderLessons принимает список ID уроков в новом порядке
func (u *Usecase) ReorderLessons(ctx context.Context, moduleID int, lessonIDs []int) error {
	orderMap := make(map[int]int)
	for i, id := range lessonIDs {
		orderMap[id] = i + 1
	}
	return u.repo.ReorderLessons(ctx, moduleID, orderMap)
}
