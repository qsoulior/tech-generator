package domain

import (
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
)

type Project struct {
	Name       string
	AuthorID   int64
	AuthorName string
	Users      []ProjectUser
}

type ProjectUser struct {
	ID   int64
	Role user_domain.Role
}
