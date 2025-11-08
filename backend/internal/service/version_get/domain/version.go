package domain

import (
	"errors"
	"time"
)

var ErrVersionNotFound = errors.New("version not found")

type Version struct {
	ID         int64
	TemplateID int64
	Number     int64
	CreatedAt  time.Time
	Data       []byte
	Variables  []Variable
}
