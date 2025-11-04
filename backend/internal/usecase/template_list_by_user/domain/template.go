package domain

import "time"

type Template struct {
	ID         int64
	Name       string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	AuthorName string
}
