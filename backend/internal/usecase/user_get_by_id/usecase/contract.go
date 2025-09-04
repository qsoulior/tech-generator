//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

type userRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}
