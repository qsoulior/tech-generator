package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrValueEmpty      = errors.New("value is empty")
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
)

type TemplateCreateIn struct {
	Name      string
	ProjectID int64
	AuthorID  int64
}

func (in TemplateCreateIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	return nil
}
