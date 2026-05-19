package domain

import user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"

type Template struct {
	AuthorID        int64
	ProjectAuthorID int64
}

type TemplateUser struct {
	ID    int64
	Name  string
	Email string
	Role  user_domain.Role
}
