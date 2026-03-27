package update

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ModuleUpdater interface {
	Update(ctx context.Context, module *domain.Module) error
}

type ModuleFinder interface {
	GetByID(ctx context.Context, id int) (*domain.Module, error)
}
