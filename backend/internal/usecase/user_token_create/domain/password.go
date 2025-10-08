package domain

import (
	"unicode/utf8"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

const (
	PasswordLengthMin = 8
	PasswordLengthMax = 72
)

var (
	ErrPasswordTooShort  = error_domain.NewBaseError("password is too short")
	ErrPasswordTooLong   = error_domain.NewBaseError("password is too long")
	ErrPasswordIncorrect = error_domain.NewBaseError("password is incorrect")
)

type Password string

func (p Password) Validate() error {
	if utf8.RuneCountInString(string(p)) < PasswordLengthMin {
		return ErrPasswordTooShort
	}

	if len(p) > PasswordLengthMax {
		return ErrPasswordTooLong
	}

	return nil
}
