package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrProjectNotFound = error_domain.NewBaseError("project not found")
	ErrProjectInvalid  = error_domain.NewBaseError("project is invalid")
	ErrValueEmpty      = errors.New("value is empty")
)

type ProjectUpdateIn struct {
	ProjectID int64
	UserID    int64
	Name      string
}

func (in ProjectUpdateIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	return nil
}
