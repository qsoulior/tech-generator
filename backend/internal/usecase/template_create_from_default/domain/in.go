package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrValueEmpty             = errors.New("value is empty")
	ErrProjectNotFound        = error_domain.NewBaseError("project not found")
	ErrProjectInvalid         = error_domain.NewBaseError("project is invalid")
	ErrSourceTemplateNotFound = error_domain.NewBaseError("source template not found")
	ErrSourceTemplateInvalid  = error_domain.NewBaseError("source template is invalid")
)

type TemplateCreateFromDefaultIn struct {
	AuthorID         int64
	ProjectID        int64
	SourceTemplateID int64
	Name             string
}

func (in TemplateCreateFromDefaultIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	return nil
}
