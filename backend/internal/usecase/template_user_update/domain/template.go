package domain

import user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"

type Template struct {
	AuthorID     int64
	RootAuthorID int64
}

type TemplateUser struct {
	ID   int64
	Role user_domain.Role
}
