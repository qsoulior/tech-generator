package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
)

var (
	ErrVersionNotFound = error_domain.NewBaseError("version not found")
	ErrVersionInvalid  = error_domain.NewBaseError("version is invalid")
)

type Version struct {
	ProjectAuthorID  int64
	TemplateAuthorID int64
	TemplateUsers    []TemplateUser
}

type TemplateUser struct {
	ID   int64
	Role user_domain.Role
}
