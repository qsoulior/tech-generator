package project_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
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
			"p.name",
			"p.author_id",
			"u.name as author_name",
			"pu.user_id",
			"pu.role",
		).
		From("project p").
		Join("usr u ON p.author_id = u.id").
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
