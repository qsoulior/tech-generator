package version_repository

import (
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

type templateVersion struct {
	ID         int64     `db:"id"`
	Number     int64     `db:"number"`
	AuthorName string    `db:"author_name"`
	CreatedAt  time.Time `db:"created_at"`
}

func (v *templateVersion) toDomain() domain.Version {
	return domain.Version{
		ID:         v.ID,
		Number:     v.Number,
		AuthorName: v.AuthorName,
		CreatedAt:  v.CreatedAt,
	}
}
