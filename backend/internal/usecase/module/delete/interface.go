package delete

import "context"

type ModuleDeleter interface {
	Delete(ctx context.Context, moduleID int) error
}
