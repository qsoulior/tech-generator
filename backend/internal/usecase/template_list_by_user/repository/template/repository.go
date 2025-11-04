package template_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ListByUserID(ctx context.Context, in domain.TemplateListByUserIn) ([]domain.Template, error) {
	op := "template - list by user id"

	subBuilder := sq.
		Select(
			"t.id",
			"t.name as template_name",
			"t.author_id as template_author_id",
			"t.created_at",
			"t.updated_at",
			"t.project_id",
			"u.name as template_author_name",
			"p.author_id as project_author_id",
			"array_agg(tu.user_id) as template_user_ids",
		).
		From("template t").
		Join("usr u ON t.author_id = u.id").
		Join("project p ON t.project_id = p.id").
		LeftJoin("template_user tu ON t.id = tu.template_id").
		GroupBy("t.id", "p.id", "u.id")

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"template_name",
			"template_author_name",
			"created_at",
			"updated_at",
		).
		From("t").
		Where(getWherePred(in.Filter)).
		OrderBy(getOrderByPred(in.Sorting)).
		Limit(uint64(in.Size)).              //nolint:gosec
		Offset(uint64((in.Page-1)*in.Size)). //nolint:gosec
		Prefix("WITH t AS (?)", subBuilder)
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

func (r *Repository) GetTotalByUserID(ctx context.Context, in domain.TemplateListByUserIn) (int64, error) {
	op := "project - get total by user id"

	subBuilder := sq.
		Select(
			"t.name as template_name",
			"t.project_id",
			"t.author_id as template_author_id",
			"p.author_id as project_author_id",
			"array_agg(tu.user_id) as template_user_ids",
		).
		From("template t").
		Join("project p ON t.project_id = p.id").
		LeftJoin("template_user tu ON t.id = tu.template_id").
		GroupBy("t.id", "p.id")

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("t").
		Where(getWherePred(in.Filter)).
		Prefix("WITH t AS (?)", subBuilder)

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

func getWherePred(filter domain.TemplateListByUserFilter) sq.Sqlizer {
	wherePred := sq.And{
		sq.Or{
			sq.Eq{"template_author_id": filter.UserID},
			sq.Eq{"project_author_id": filter.UserID},
			sq.Expr("? = ANY(template_user_ids)", filter.UserID),
		},
		sq.Eq{"project_id": filter.ProjectID},
	}

	if filter.TemplateName != nil {
		wherePred = append(wherePred, sq.ILike{"template_name": fmt.Sprintf("%%%s%%", *filter.TemplateName)})
	}

	return wherePred
}

func getOrderByPred(sorting *sorting_domain.Sorting) string {
	if sorting == nil {
		return "id DESC"
	}

	if _, ok := sortingAttributes[sorting.Attribute]; !ok {
		return "id DESC"
	}

	return fmt.Sprintf("%s %s", sorting.Attribute, sorting.Direction)
}
