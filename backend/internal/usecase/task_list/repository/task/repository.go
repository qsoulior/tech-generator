package task_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context, in domain.TaskListIn) ([]domain.Task, error) {
	op := "task - list"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"t.id",
			"v.number as version_number",
			"t.status",
			"u.name as creator_name",
			"t.created_at",
			"t.updated_at",
		).
		From("task t").
		Join("template_version v ON t.version_id = v.id").
		Join("usr u ON t.creator_id = u.id").
		Where(getWherePred(in.Filter)).
		OrderBy(getOrderByPred(in.Sorting)).
		Limit(uint64(in.Size)).                 //nolint:gosec
		Offset(uint64((in.Page - 1) * in.Size)) //nolint:gosec

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []task
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	projects := lo.Map(dtos, func(dto task, _ int) domain.Task { return dto.toDomain() })
	return projects, nil
}

func (r *Repository) GetTotal(ctx context.Context, in domain.TaskListIn) (int64, error) {
	op := "task - get total"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("task t").
		Join("template_version v ON t.version_id = v.id").
		Join("usr u ON t.creator_id = u.id").
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

func getWherePred(filter domain.TaskListFilter) sq.Sqlizer {
	wherePred := sq.And{
		sq.Eq{"v.template_id": filter.TemplateID},
	}

	if filter.CreatorID != nil {
		wherePred = append(wherePred, sq.Eq{"creator_id": *filter.CreatorID})
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
