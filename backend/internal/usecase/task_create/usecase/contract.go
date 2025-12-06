//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

type versionRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Version, error)
}

type taskRepository interface {
	Insert(ctx context.Context, in domain.TaskCreateIn) (int64, error)
}

type publisher interface {
	PublishTaskCreated(ctx context.Context, id int64) error
}
