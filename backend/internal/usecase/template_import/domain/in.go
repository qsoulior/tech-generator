package domain

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
)

var (
	ErrValueEmpty      = errors.New("value is empty")
	ErrValueInvalid    = errors.New("value is invalid")
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
)

var slugRegexp = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

const (
	slugMaxLen  = 100
	titleMaxLen = 255
)

type Constraint struct {
	Name       string
	Expression string
	IsActive   bool
}

type Variable struct {
	Name        string
	Title       string
	Type        variable_domain.Type
	Expression  *string
	IsInput     bool
	Constraints []Constraint
}

type Version struct {
	Data      []byte
	Variables []Variable
}

type TemplateImportIn struct {
	AuthorID  int64
	ProjectID int64
	Name      string
	Version   *Version
}

func (in TemplateImportIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	if in.Version == nil {
		return nil
	}

	for i, v := range in.Version.Variables {
		if !v.Type.Valid() {
			return error_domain.NewValidationError(fmt.Sprintf("version.variables.%d.type", i), ErrValueInvalid)
		}

		if v.Name == "" || len(v.Name) > slugMaxLen || !slugRegexp.MatchString(v.Name) {
			return error_domain.NewValidationError(fmt.Sprintf("version.variables.%d.name", i), ErrValueInvalid)
		}

		if v.Title == "" {
			return error_domain.NewValidationError(fmt.Sprintf("version.variables.%d.title", i), ErrValueEmpty)
		}

		if utf8.RuneCountInString(v.Title) > titleMaxLen {
			return error_domain.NewValidationError(fmt.Sprintf("version.variables.%d.title", i), ErrValueInvalid)
		}
	}

	return nil
}
