package domain

import (
	"errors"

	base_domain "github.com/qsoulior/tech-generator/backend/internal/domain"
)

var (
	ErrEmptyValue = errors.New("value is empty")
	ErrUserExists = base_domain.NewError("user with name or email exists")
)

type UserCreateIn struct {
	Name     string
	Email    string
	Password Password
}

func (in UserCreateIn) Validate() error {
	if in.Email == "" {
		return NewValidationError("email", ErrEmptyValue)
	}

	if in.Name == "" {
		return NewValidationError("name", ErrEmptyValue)
	}

	if err := in.Password.Validate(); err != nil {
		return NewValidationError("password", err)
	}

	return nil
}
