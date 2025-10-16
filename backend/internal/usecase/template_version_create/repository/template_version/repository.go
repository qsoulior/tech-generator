package template_version_repository

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

func (r *Repository) Create(ctx context.Context, templateVersion domain.TemplateVersion) (int64, error) {
	op := "template version - create"

	numberExpr := sq.Expr("(SELECT COALESCE(MAX(number), 0) FROM template_version WHERE template_id = ?) + 1", templateVersion.TemplateID)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("template_version").
		Columns("number", "template_id", "author_id", "data").
		Values(
			numberExpr,
			templateVersion.TemplateID,
			templateVersion.AuthorID,
			templateVersion.Data,
		).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var id int64
	err = r.trGetter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query %q: %w", op, err)
	}

	return id, nil
}
