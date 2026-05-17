//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
)

type projectRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Project, error)
}
