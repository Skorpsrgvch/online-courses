package create

import (
	"context"

	"github.com/Skorpsrgvch/online-courses/internal/domain"
)

type ModuleSaver interface {
	Save(ctx context.Context, module *domain.Module) error
}
