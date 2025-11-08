package domain

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
)

type Template struct {
	LastVersionID   *int64
	AuthorID        int64
	ProjectAuthorID int64
	Users           []TemplateUser
}

type TemplateUser struct {
	ID   int64
	Role user_domain.Role
}
