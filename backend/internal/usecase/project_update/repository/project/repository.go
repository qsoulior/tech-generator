package project_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/domain"
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
		Select(
			"p.author_id",
			"pu.user_id",
			"pu.role",
		).
		From("project p").
		LeftJoin("project_user pu ON p.id = pu.project_id").
		Where(sq.Eq{"p.id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos projects
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return dtos.toDomain(), nil
}

func (r *Repository) UpdateByID(ctx context.Context, id int64, name string) error {
	op := "project - update by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("project").
		Set("name", name).
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
