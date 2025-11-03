package folder_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Folder, error) {
	op := "folder - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"author_id",
			"root_author_id",
		).
		From("folder").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dto folder
	err = r.db.GetContext(ctx, &dto, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return dto.toDomain(), nil
}
