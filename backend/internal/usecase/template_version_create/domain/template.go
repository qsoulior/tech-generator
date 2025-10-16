package domain

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
)

type Template struct {
	AuthorID     int64
	RootAuthorID int64
	Users        []TemplateUser
}

type TemplateUser struct {
	ID   int64
	Role user_domain.Role
}

type TemplateToUpdate struct {
	ID            int64
	LastVersionID int64
}

type TemplateVersion struct {
	TemplateID int64
	AuthorID   int64
	Data       []byte
}
