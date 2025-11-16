package task_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	op := "task - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"version_id",
			"payload",
		).
		From("task").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dto task
	err = r.db.GetContext(ctx, &dto, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return dto.toDomain(), nil
}

func (r *Repository) UpdateByID(ctx context.Context, task domain.TaskUpdate) error {
	op := "task - update by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("task").
		SetMap(map[string]any{
			"status":     task.Status,
			"result":     task.Result,
			"error":      (*taskError)(task.Error),
			"updated_at": sq.Expr("now() AT TIME ZONE 'utc'"),
		}).
		Where(sq.Eq{"id": task.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query %q: %w", op, err)
	}

	return nil
}
