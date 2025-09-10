//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import (
	"context"
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type userRepository interface {
	GetByName(ctx context.Context, name string) (*domain.User, error)
}

type passwordVerifier interface {
	Verify(passwordHash []byte, password domain.Password) error
}

type tokenBuilder interface {
	Build(userID int64, tokenExpiration time.Duration) (domain.UserCreateTokenOut, error)
}
