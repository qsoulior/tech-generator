package domain

import user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"

type Project struct {
	AuthorID int64
	Users    []ProjectUser
}

type ProjectUser struct {
	ID   int64
	Role user_domain.Role
}
