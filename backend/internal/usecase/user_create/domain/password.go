package domain

import (
	"errors"
	"slices"
	"unicode"
	"unicode/utf8"
)

const (
	PasswordLengthMin = 8
	PasswordLengthMax = 72
)

var (
	ErrPasswordTooShort           = errors.New("password is too short")
	ErrPasswordTooLong            = errors.New("password is too long")
	ErrPasswordNoDigits           = errors.New("no digits in password")
	ErrPasswordNoLettersUppercase = errors.New("no uppercase letters in password")
	ErrPasswordNoLettersLowercase = errors.New("no lowercase letters in password")
	ErrPasswordNoSpecial          = errors.New("no special character in password")
)

type Password string

func (p Password) Validate() error {
	cnt := utf8.RuneCountInString(string(p))

	if cnt < PasswordLengthMin {
		return ErrPasswordTooShort
	}

	if cnt > PasswordLengthMax {
		return ErrPasswordTooLong
	}

	runes := []rune(p)

	hasDigit := slices.ContainsFunc(runes, func(r rune) bool { return unicode.IsDigit(r) })
	if !hasDigit {
		return ErrPasswordNoDigits
	}

	hasLetterUppercase := slices.ContainsFunc(runes, func(r rune) bool { return unicode.IsUpper(r) })
	if !hasLetterUppercase {
		return ErrPasswordNoLettersUppercase
	}

	hasLetterLowercase := slices.ContainsFunc(runes, func(r rune) bool { return unicode.IsLower(r) })
	if !hasLetterLowercase {
		return ErrPasswordNoLettersLowercase
	}

	hasSpecial := slices.ContainsFunc(runes, func(r rune) bool { return unicode.IsPunct(r) || unicode.IsSymbol(r) })
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}
