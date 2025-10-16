package variable_constraint_repository

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

func (r *Repository) Create(ctx context.Context, constraints []domain.VariableConstraintToCreate) error {
	op := "variable constraint - create"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("variable_constraint").
		Columns("variable_id", "name", "expression", "is_active")

	for _, c := range constraints {
		builder = builder.Values(c.VariableID, c.Name, c.Expression, c.IsActive)
	}

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
