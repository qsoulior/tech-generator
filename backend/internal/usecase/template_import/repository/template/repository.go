package template_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, template domain.Template) (int64, error) {
	op := "template - create"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("template").
		Columns("name", "is_default", "project_id", "author_id").
		Values(
			template.Name,
			template.IsDefault,
			template.ProjectID,
			template.AuthorID,
		).
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
