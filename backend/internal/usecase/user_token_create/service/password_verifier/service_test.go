package password_verifier

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

func TestService_Verify_Success(t *testing.T) {
	service := New()

	password := domain.Password("aB12345-")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	err = service.Verify(passwordHash, password)
	require.NoError(t, err)
}

func TestService_Verify_Error(t *testing.T) {
	service := New()

	password := domain.Password(strings.Repeat("-", 73))
	err := service.Verify([]byte{1}, password)
	require.Error(t, err)
}
