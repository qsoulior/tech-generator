package user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int64) ([]int64, error) {
	op := "usr - get by ids"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id").
		From("usr").
		Where(sq.Eq{"id": ids})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var userIDs []int64
	err = r.db.SelectContext(ctx, &userIDs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return userIDs, nil
}
