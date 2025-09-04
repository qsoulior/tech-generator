//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

type userRepository interface {
	Create(ctx context.Context, name, email string, password []byte) error
	ExistsByNameOrEmail(ctx context.Context, name, email string) (bool, error)
}

type passwordHasher interface {
	Hash(password domain.Password) ([]byte, error)
}
