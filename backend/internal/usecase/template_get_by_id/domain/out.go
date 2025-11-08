package domain

import "time"

type TemplateGetByIDOut struct {
	VersionID     int64
	VersionNumber int64
	CreatedAt     time.Time
	Data          []byte
	Variables     []Variable
}
