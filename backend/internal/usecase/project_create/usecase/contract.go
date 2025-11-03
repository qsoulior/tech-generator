//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

type projectRepository interface {
	Create(ctx context.Context, project domain.Project) error
}
