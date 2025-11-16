//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type taskRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Task, error)
	UpdateByID(ctx context.Context, task domain.TaskUpdate) error
}

type versionGetService interface {
	Handle(ctx context.Context, versionID int64) (*version_get_domain.Version, error)
}

type variableProcessService interface {
	Handle(ctx context.Context, in domain.VariableProcessIn) (map[string]any, error)
}

type dataProcessService interface {
	Handle(ctx context.Context, in domain.DataProcessIn) ([]byte, error)
}
