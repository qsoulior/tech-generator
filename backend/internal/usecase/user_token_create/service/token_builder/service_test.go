package token_builder

import (
	"crypto/ed25519"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestService_Build(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	wantUserID := gofakeit.Int64()
	tokenExpiration := 7 * 24 * time.Hour

	service := New(privateKey)
	out, err := service.Build(wantUserID, tokenExpiration)
	require.NoError(t, err)

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer("tech-generator"),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	// validate token
	gotToken, err := parser.ParseWithClaims(out.Token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) { return publicKey, nil })
	require.NoError(t, err)
	require.True(t, gotToken.Valid)

	// validate subject
	subject, err := gotToken.Claims.GetSubject()
	require.NoError(t, err)

	gotUserID, err := strconv.ParseInt(subject, 10, 64)
	require.NoError(t, err)
	require.Equal(t, wantUserID, gotUserID)

	// validate expiresAt
	expiresAt, err := gotToken.Claims.GetExpirationTime()
	require.NoError(t, err)
	require.Equal(t, expiresAt.UTC(), out.ExpiresAt.Truncate(1*time.Second))
}
