package template_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ListDefault(ctx context.Context, in domain.TemplateListDefaultIn) ([]domain.Template, error) {
	op := "template - list default"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"name",
			"created_at",
			"updated_at",
		).
		From("template").
		Where(getWherePred(in.Filter)).
		OrderBy(getOrderByPred(in.Sorting)).
		Limit(uint64(in.Size)).                 //nolint:gosec
		Offset(uint64((in.Page - 1) * in.Size)) //nolint:gosec

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []template
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	templates := lo.Map(dtos, func(dto template, _ int) domain.Template { return dto.toDomain() })
	return templates, nil
}

func (r *Repository) GetTotalDefault(ctx context.Context, in domain.TemplateListDefaultIn) (int64, error) {
	op := "template - get total default"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("template").
		Where(getWherePred(in.Filter))

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var count int64
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query %q: %w", op, err)
	}

	return count, nil
}

func getWherePred(filter domain.TemplateListDefaultFilter) sq.Sqlizer {
	wherePred := sq.And{
		sq.Eq{"is_default": true},
	}

	if filter.TemplateName != nil {
		wherePred = append(wherePred, sq.ILike{"name": fmt.Sprintf("%%%s%%", *filter.TemplateName)})
	}

	return wherePred
}

func getOrderByPred(sorting *sorting_domain.Sorting) string {
	const defaultOrder = "COALESCE(updated_at, created_at) DESC, id DESC"

	if sorting == nil {
		return defaultOrder
	}

	if _, ok := sortingAttributes[sorting.Attribute]; !ok {
		return defaultOrder
	}

	return fmt.Sprintf("%s %s", sorting.Attribute, sorting.Direction)
}
