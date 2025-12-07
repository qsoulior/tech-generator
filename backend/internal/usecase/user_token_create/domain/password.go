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
	ErrPasswordTooShort  = error_domain.NewBaseError("Пароль слишком короткий")
	ErrPasswordTooLong   = error_domain.NewBaseError("Пароль слишком длинный")
	ErrPasswordIncorrect = error_domain.NewBaseError("Имя пользователя или пароль неверные")
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
