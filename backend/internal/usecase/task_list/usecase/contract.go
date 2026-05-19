//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type templateRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Template, error)
}

type taskRepository interface {
	List(ctx context.Context, in domain.TaskListIn) ([]domain.Task, error)
	GetTotal(ctx context.Context, in domain.TaskListIn) (int64, error)
}
