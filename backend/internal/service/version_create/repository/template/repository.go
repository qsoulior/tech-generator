package template_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
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
