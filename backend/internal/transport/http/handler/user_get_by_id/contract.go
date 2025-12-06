//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package user_get_by_id_handler

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

type usecase interface {
	Handle(ctx context.Context, id int64) (*domain.User, error)
}
