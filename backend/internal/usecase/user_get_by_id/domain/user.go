package domain

import (
	"time"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

var ErrUserNotFound = error_domain.NewBaseError("user not found")

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}
