package constraint_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ListByVariableIDs(ctx context.Context, variableIDs []int64) ([]domain.Constraint, error) {
	op := "constraint - list by variable ids"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"variable_id",
			"name",
			"expression",
			"is_active",
		).
		From("variable_constraint").
		Where(sq.Eq{"variable_id": variableIDs}).
		OrderBy("id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []сonstraint
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	constraints := lo.Map(dtos, func(dto сonstraint, _ int) domain.Constraint { return dto.toDomain() })
	return constraints, nil
}
