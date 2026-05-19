//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

type userRepository interface {
	List(ctx context.Context, in domain.UserListIn) ([]domain.User, error)
	GetTotal(ctx context.Context, in domain.UserListIn) (int64, error)
}
