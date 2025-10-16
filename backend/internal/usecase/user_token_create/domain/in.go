package domain

import error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"

type UserCreateTokenIn struct {
	Name     string
	Password Password
}

var ErrNameEmpty = error_domain.NewBaseError("name is empty")

func (in UserCreateTokenIn) Validate() error {
	if in.Name == "" {
		return ErrNameEmpty
	}

	return in.Password.Validate()
}
