package project_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ListByUserID(ctx context.Context, in domain.ProjectListByUserIn) ([]domain.Project, error) {
	op := "project - list by user id"

	subBuilder := sq.
		Select(
			"p.id",
			"p.name as project_name",
			"p.author_id",
			"u.name as author_name",
			"array_agg(pu.user_id) as project_user_ids",
		).
		From("project p").
		Join("usr u ON p.author_id = u.id").
		LeftJoin("project_user pu ON p.id = pu.project_id").
		GroupBy("p.id", "u.name")

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"project_name",
			"author_name",
		).
		From("p").
		Where(getWherePred(in.Filter)).
		OrderBy(getOrderByPred(in.Sorting)).
		Limit(uint64(in.Size)).              //nolint:gosec
		Offset(uint64((in.Page-1)*in.Size)). //nolint:gosec
		Prefix("WITH p AS (?)", subBuilder)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []project
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	projects := lo.Map(dtos, func(dto project, _ int) domain.Project { return dto.toDomain() })
	return projects, nil
}

func (r *Repository) GetTotalByUserID(ctx context.Context, in domain.ProjectListByUserIn) (int64, error) {
	op := "project - get total by user id"

	subBuilder := sq.
		Select(
			"p.name as project_name",
			"p.author_id",
			"array_agg(pu.user_id) as project_user_ids",
		).
		From("project p").
		LeftJoin("project_user pu ON p.id = pu.project_id").
		GroupBy("p.id")

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("p").
		Where(getWherePred(in.Filter)).
		Prefix("WITH p AS (?)", subBuilder)

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

func getWherePred(filter domain.ProjectListByUserFilter) sq.Sqlizer {
	wherePred := sq.And{}

	wherePred = append(wherePred, sq.Or{
		sq.Eq{"author_id": filter.UserID},
		sq.Expr("? = ANY(project_user_ids)", filter.UserID),
	})

	if filter.ProjectName != nil {
		wherePred = append(wherePred, sq.ILike{"project_name": fmt.Sprintf("%%%s%%", *filter.ProjectName)})
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
