package token_parser

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse/domain"
)

func TestService_Parse(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(24 * time.Hour)

	claims := jwt.RegisteredClaims{
		Issuer:    "tech-generator",
		Subject:   "123",
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	service := New(publicKey)
	gotUser, err := service.Parse(tokenString)
	require.NoError(t, err)

	wantUser := domain.User{ID: 123}
	require.NotNil(t, gotUser)
	require.Equal(t, wantUser, *gotUser)
}

func TestService_Parse_Error(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	service := New(publicKey)

	tests := []struct {
		name       string
		setupToken func() (string, error)
		want       string
	}{
		{
			name: "SignatureInvalid",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()
				expiresAt := issuedAt.Add(24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Issuer:    "tech-generator",
					Subject:   "123",
					IssuedAt:  jwt.NewNumericDate(issuedAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				_, privateKey, err := ed25519.GenerateKey(nil)
				if err != nil {
					return "", err
				}

				return token.SignedString(privateKey)
			},
			want: "signature",
		},
		{
			name: "TokenBeforeIssued",
			setupToken: func() (string, error) {
				expiresAt := time.Now().UTC().Add(24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Issuer:    "tech-generator",
					Subject:   "123",
					IssuedAt:  jwt.NewNumericDate(expiresAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "issued",
		},
		{
			name: "TokenExpired",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()
				expiresAt := issuedAt.Add(-24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Issuer:    "tech-generator",
					Subject:   "123",
					IssuedAt:  jwt.NewNumericDate(issuedAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "expired",
		},
		{
			name: "ExpiresAtRequired",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()

				claims := jwt.RegisteredClaims{
					Issuer:   "tech-generator",
					Subject:  "123",
					IssuedAt: jwt.NewNumericDate(issuedAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "exp",
		},
		{
			name: "IssuerEmpty",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()
				expiresAt := issuedAt.Add(24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Subject:   "123",
					IssuedAt:  jwt.NewNumericDate(issuedAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "iss",
		},
		{
			name: "IssuerInvalid",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()
				expiresAt := issuedAt.Add(24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Issuer:    "invalid",
					Subject:   "123",
					IssuedAt:  jwt.NewNumericDate(issuedAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "iss",
		},
		{
			name: "SubjectInvalid",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()
				expiresAt := issuedAt.Add(24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Issuer:    "tech-generator",
					Subject:   "invalid",
					IssuedAt:  jwt.NewNumericDate(issuedAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "parse int",
		},
		{
			name: "SubjectEmpty",
			setupToken: func() (string, error) {
				issuedAt := time.Now().UTC()
				expiresAt := issuedAt.Add(24 * time.Hour)

				claims := jwt.RegisteredClaims{
					Issuer:    "tech-generator",
					IssuedAt:  jwt.NewNumericDate(issuedAt),
					ExpiresAt: jwt.NewNumericDate(expiresAt),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
				return token.SignedString(privateKey)
			},
			want: "parse int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := tt.setupToken()
			require.NoError(t, err)

			_, err = service.Parse(tokenString)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
