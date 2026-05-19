package domain

import (
	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
)

var (
	ErrTemplateNotFound = error_domain.NewBaseError("template not found")
	ErrTemplateInvalid  = error_domain.NewBaseError("template is invalid")
	ErrRoleInvalid      = error_domain.NewBaseError("role is invalid for template")
)

type TemplateUserUpdateIn struct {
	UserID     int64
	TemplateID int64
	Users      []TemplateUser
}

func (in TemplateUserUpdateIn) Validate() error {
	for _, u := range in.Users {
		if u.Role != user_domain.RoleRead && u.Role != user_domain.RoleWrite {
			return ErrRoleInvalid
		}
	}
	return nil
}
