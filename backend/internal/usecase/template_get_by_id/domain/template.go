package domain

import (
	"time"

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

type TemplateVersion struct {
	ID        int64
	Number    int64
	CreatedAt time.Time
	Data      []byte
}
