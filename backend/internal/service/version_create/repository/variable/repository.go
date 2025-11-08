package variable_repository

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

func (r *Repository) Create(ctx context.Context, variables []domain.VariableToCreate) ([]int64, error) {
	op := "variable - create"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("variable").
		Columns("version_id", "name", "type", "expression").
		Suffix("RETURNING id")

	for _, v := range variables {
		builder = builder.Values(v.VersionID, v.Name, v.Type, v.Expression)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var ids []int64
	err = r.trGetter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &ids, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return ids, nil
}
