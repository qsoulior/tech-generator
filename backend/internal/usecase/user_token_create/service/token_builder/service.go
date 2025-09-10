package token_builder

import (
	"crypto/ed25519"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type Service struct {
	privateKey ed25519.PrivateKey
}

func New(privateKey ed25519.PrivateKey) *Service {
	return &Service{
		privateKey: privateKey,
	}
}

func (s *Service) Build(userID int64, tokenExpiration time.Duration) (domain.UserCreateTokenOut, error) {
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(tokenExpiration)

	claims := jwt.RegisteredClaims{
		Issuer:    "tech-generator",
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return domain.UserCreateTokenOut{}, fmt.Errorf("token - signed string")
	}

	out := domain.UserCreateTokenOut{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	return out, nil
}
