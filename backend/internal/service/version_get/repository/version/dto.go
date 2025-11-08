package version_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type version struct {
	ID        int64     `db:"id"`
	Number    int64     `db:"number"`
	CreatedAt time.Time `db:"created_at"`
	Data      []byte    `db:"data"`
}

func (v *version) toDomain() *domain.Version {
	return &domain.Version{
		ID:        v.ID,
		Number:    v.Number,
		CreatedAt: v.CreatedAt,
		Data:      v.Data,
	}
}
