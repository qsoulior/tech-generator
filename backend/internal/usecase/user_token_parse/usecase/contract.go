//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package usecase

import "github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"

type tokenParser interface {
	Parse(tokenString string) (*domain.User, error)
}
