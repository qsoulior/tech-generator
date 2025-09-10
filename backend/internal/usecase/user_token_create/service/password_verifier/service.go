package password_verifier

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Verify(passwordHash []byte, password domain.Password) error {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		return fmt.Errorf("bcrypt - compare hash and password: %w", err)
	}

	return nil
}
