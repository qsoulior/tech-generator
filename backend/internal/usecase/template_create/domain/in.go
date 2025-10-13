package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrEmptyValue     = errors.New("value is empty")
	ErrFolderNotFound = error_domain.NewBaseError("folder not found")
	ErrFolderInvalid  = error_domain.NewBaseError("folder is invalid")
)

type TemplateCreateIn struct {
	Name     string
	FolderID int64
	AuthorID int64
}

func (in TemplateCreateIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrEmptyValue)
	}

	return nil
}
