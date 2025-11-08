package domain

import (
	"time"
)

type Template struct {
	AuthorID        int64
	ProjectAuthorID int64
}

type TemplateVersion struct {
	ID         int64
	Number     int64
	AuthorName string
	CreatedAt  time.Time
}
