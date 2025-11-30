//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package auth_middleware

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"
)

type usecase interface {
	Handle(_ context.Context, token string) (*domain.User, error)
}
