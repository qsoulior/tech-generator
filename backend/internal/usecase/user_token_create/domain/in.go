package domain

import (
	base_domain "github.com/qsoulior/tech-generator/backend/internal/domain"
)

type UserCreateTokenIn struct {
	Name     string
	Password Password
}

var ErrNameEmptyValue = base_domain.NewError("name is empty")

func (in UserCreateTokenIn) Validate() error {
	if in.Name == "" {
		return ErrNameEmptyValue
	}

	return in.Password.Validate()
}
