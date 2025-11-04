//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

type projectRepository interface {
	ListByUserID(ctx context.Context, in domain.ProjectListByUserIn) ([]domain.Project, error)
	GetTotalByUserID(ctx context.Context, in domain.ProjectListByUserIn) (int64, error)
}
