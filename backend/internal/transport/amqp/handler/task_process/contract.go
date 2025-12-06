//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package task_process_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TaskProcessIn) error
}
