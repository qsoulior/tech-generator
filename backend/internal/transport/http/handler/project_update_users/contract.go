//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package project_update_users_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.ProjectUserUpdateIn) error
}
