//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package project_users_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.ProjectUserListIn) ([]domain.ProjectUser, error)
}
