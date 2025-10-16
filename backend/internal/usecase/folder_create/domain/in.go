package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrValueEmpty     = errors.New("value is empty")
	ErrParentNotFound = error_domain.NewBaseError("parent folder not found")
	ErrParentInvalid  = error_domain.NewBaseError("parent folder is invalid")
)

type FolderCreateIn struct {
	ParentID *int64
	Name     string
	AuthorID int64
}

func (in FolderCreateIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	return nil
}
