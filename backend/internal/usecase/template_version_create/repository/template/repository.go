package template_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

type Repository struct {
	db       *sqlx.DB
	trGetter *trmsqlx.CtxGetter
}

func New(db *sqlx.DB, trGetter *trmsqlx.CtxGetter) *Repository {
	return &Repository{
		db:       db,
		trGetter: trGetter,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Template, error) {
	op := "template - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"t.author_id",
			"t.root_author_id",
			"tu.user_id",
			"tu.role",
		).
		From("template t").
		LeftJoin("template_user tu ON t.id = tu.template_id").
		Where(sq.Eq{"t.id": id, "t.is_default": false})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var templates templates
	err = r.db.SelectContext(ctx, &templates, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return templates.toDomain(), nil
}

func (r *Repository) UpdateByID(ctx context.Context, template domain.TemplateToUpdate) error {
	op := "template - update by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("template").
		Set("last_version_id", template.LastVersionID).
		Set("updated_at", sq.Expr("now() AT TIME ZONE 'utc'")).
		Where(sq.Eq{"id": template.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	_, err = r.trGetter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query %q: %w", op, err)
	}

	return nil
}
