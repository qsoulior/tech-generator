package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
)

var (
	ErrTemplateNotFound = error_domain.NewBaseError("template not found")
	ErrTemplateInvalid  = error_domain.NewBaseError("template is invalid")
)

type Template struct {
	AuthorID        int64
	ProjectAuthorID int64
	Users           []TemplateUser
}

type TemplateUser struct {
	ID   int64
	Role user_domain.Role
}
