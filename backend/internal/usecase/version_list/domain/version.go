package domain

import "time"

type Version struct {
	ID         int64
	Number     int64
	AuthorName string
	CreatedAt  time.Time
}
