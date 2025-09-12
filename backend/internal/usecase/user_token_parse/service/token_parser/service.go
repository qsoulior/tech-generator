package token_parser

import (
	"crypto/ed25519"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"
)

type Service struct {
	publicKey ed25519.PublicKey
	parser    *jwt.Parser
}

func New(publicKey ed25519.PublicKey) *Service {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer("tech-generator"),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	return &Service{
		publicKey: publicKey,
		parser:    parser,
	}
}

func (s *Service) Parse(tokenString string) (*domain.User, error) {
	token, err := s.parser.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) { return s.publicKey, nil })
	if err != nil {
		return nil, fmt.Errorf("jwt parser - parse with claims: %w", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("token claims - get subject: %w", err)
	}

	userID, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv - parse int %q: %w", subject, err)
	}

	user := domain.User{ID: userID}
	return &user, nil
}
