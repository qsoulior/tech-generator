package domain

import (
	"errors"
	"fmt"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var (
	ErrValueInvalid       = errors.New("value is invalid")
	ErrTemplateNotFound   = error_domain.NewBaseError("template not found")
	ErrTemplateInvalid    = error_domain.NewBaseError("template is invalid")
	ErrVariableIDsInvalid = errors.New("variable ids length is invalid")
)

type TemplateVersionCreateIn struct {
	AuthorID   int64
	TemplateID int64
	Data       []byte
	Variables  []Variable
}

func (in TemplateVersionCreateIn) Validate() error {
	for i, v := range in.Variables {
		if !v.Type.Valid() {
			field := fmt.Sprintf("variables.%d.type", i)
			return error_domain.NewValidationError(field, ErrValueInvalid)
		}
	}

	return nil
}
