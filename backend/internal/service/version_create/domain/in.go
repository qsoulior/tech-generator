package domain

import (
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrValueInvalid       = errors.New("value is invalid")
	ErrVariableIDsInvalid = errors.New("variable ids length is invalid")
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
			field := fmt.Sprintf("variables.%d.type", i)
			return error_domain.NewValidationError(field, ErrValueInvalid)
		}
	}

	return nil
}
