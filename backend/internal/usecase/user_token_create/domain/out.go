package domain

import "time"

type UserCreateTokenOut struct {
	Token     string
	ExpiresAt time.Time
}
