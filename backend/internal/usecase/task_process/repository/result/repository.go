package result_repository

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

func (r *Repository) Insert(ctx context.Context, data []byte) (int64, error) {
	op := "result - insert"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("result").
		Columns("data").
		Values(data).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var id int64
	err = r.db.GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query %q: %w", op, err)
	}

	return id, nil
}
