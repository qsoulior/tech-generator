package domain

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrValueInvalid       = errors.New("value is invalid")
	ErrValueEmpty         = errors.New("value is empty")
	ErrVariableIDsInvalid = errors.New("variable ids length is invalid")
)

var slugRegexp = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

const (
	slugMaxLen  = 100
	titleMaxLen = 255
)

type VersionCreateIn struct {
	AuthorID   int64
	TemplateID int64
	Data       []byte
	Variables  []Variable
}

func (in VersionCreateIn) Validate() error {
	for i, v := range in.Variables {
		if !v.Type.Valid() {
			return error_domain.NewValidationError(fmt.Sprintf("variables.%d.type", i), ErrValueInvalid)
		}

		if v.Name == "" || len(v.Name) > slugMaxLen || !slugRegexp.MatchString(v.Name) {
			return error_domain.NewValidationError(fmt.Sprintf("variables.%d.name", i), ErrValueInvalid)
		}

		if v.Title == "" {
			return error_domain.NewValidationError(fmt.Sprintf("variables.%d.title", i), ErrValueEmpty)
		}

		if utf8.RuneCountInString(v.Title) > titleMaxLen {
			return error_domain.NewValidationError(fmt.Sprintf("variables.%d.title", i), ErrValueInvalid)
		}
	}

	return nil
}
