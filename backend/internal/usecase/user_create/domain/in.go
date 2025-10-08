package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrEmptyValue = errors.New("value is empty")
	ErrUserExists = error_domain.NewBaseError("user with name or email exists")
)

type UserCreateIn struct {
	Name     string
	Email    string
	Password Password
}

func (in UserCreateIn) Validate() error {
	if in.Email == "" {
		return error_domain.NewValidationError("email", ErrEmptyValue)
	}

	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrEmptyValue)
	}

	if err := in.Password.Validate(); err != nil {
		return error_domain.NewValidationError("password", err)
	}

	return nil
}
