//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package template_users_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.TemplateUserListIn) ([]domain.TemplateUser, error)
}
