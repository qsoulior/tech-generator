package password_hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Hash(password domain.Password) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt - generate from password: %w", err)
	}

	return hash, nil
}
