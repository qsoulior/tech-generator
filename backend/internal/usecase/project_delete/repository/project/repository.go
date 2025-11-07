package project_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	op := "project - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("author_id").
		From("project").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var project project
	err = r.db.GetContext(ctx, &project, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return project.toDomain(), nil
}

func (r *Repository) DeleteByID(ctx context.Context, id int64) error {
	op := "project - delete by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("project").
		Where(sq.Eq{"id": id})

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
