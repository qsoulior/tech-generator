//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package user_token_create_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type usecase interface {
	Handle(ctx context.Context, in domain.UserCreateTokenIn) (domain.UserCreateTokenOut, error)
}
