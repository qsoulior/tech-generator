package domain

import (
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
)

var (
	ErrValueEmpty      = errors.New("value is empty")
	ErrValueInvalid    = errors.New("value is invalid")
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
)

type Constraint struct {
	Name       string
	Expression string
	IsActive   bool
}

type Variable struct {
	Name        string
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
			field := fmt.Sprintf("version.variables.%d.type", i)
			return error_domain.NewValidationError(field, ErrValueInvalid)
		}
	}

	return nil
}
