package password_hasher

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

func TestService_Hash_Success(t *testing.T) {
	service := New()

	password := domain.Password("aB12345-")
	got, err := service.Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, got)
}

func TestService_Hash_Error(t *testing.T) {
	service := New()

	password := domain.Password(strings.Repeat("-", 73))
	_, err := service.Hash(password)
	require.Error(t, err)
}
