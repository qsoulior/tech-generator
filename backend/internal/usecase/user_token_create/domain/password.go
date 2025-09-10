package domain

import (
	"unicode/utf8"

	base_domain "github.com/qsoulior/tech-generator/backend/internal/domain"
)

const (
	PasswordLengthMin = 8
	PasswordLengthMax = 72
)

var (
	ErrPasswordTooShort  = base_domain.NewError("password is too short")
	ErrPasswordTooLong   = base_domain.NewError("password is too long")
	ErrPasswordIncorrect = base_domain.NewError("password is incorrect")
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
