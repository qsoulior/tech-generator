package domain

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
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
