package domain

import (
	"errors"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var ErrValueEmpty = errors.New("value is empty")

type ProjectCreateIn struct {
	Name     string
	AuthorID int64
}

func (in ProjectCreateIn) Validate() error {
	if in.Name == "" {
		return error_domain.NewValidationError("name", ErrValueEmpty)
	}

	return nil
}
