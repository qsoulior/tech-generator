//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	domain "github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
)

type taskRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Task, error)
}

type versionRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Version, error)
}

type resultRepository interface {
	GetDataByID(ctx context.Context, id int64) ([]byte, error)
}
