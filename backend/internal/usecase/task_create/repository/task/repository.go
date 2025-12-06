package task_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Insert(ctx context.Context, in domain.TaskCreateIn) (int64, error) {
	op := "task - insert"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("task").
		Columns("version_id", "creator_id", "payload").
		Values(in.VersionID, in.CreatorID, payload(in.Payload)).
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
