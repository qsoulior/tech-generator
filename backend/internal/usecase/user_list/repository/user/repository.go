package user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context, in domain.UserListIn) ([]domain.User, error) {
	op := "usr - list"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"name",
			"email",
		).
		From("usr").
		Where(getWherePred(in.Filter)).
		OrderBy("name ASC").
		Limit(uint64(in.Size)).                 //nolint:gosec
		Offset(uint64((in.Page - 1) * in.Size)) //nolint:gosec

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []user
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	users := lo.Map(dtos, func(dto user, _ int) domain.User { return dto.toDomain() })
	return users, nil
}

func (r *Repository) GetTotal(ctx context.Context, in domain.UserListIn) (int64, error) {
	op := "usr - get total"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("usr").
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

func getWherePred(filter domain.UserListFilter) sq.Sqlizer {
	wherePred := sq.And{
		sq.NotEq{"id": filter.ExcludeUserID},
	}

	if filter.UserName != nil {
		wherePred = append(wherePred, sq.ILike{"name": fmt.Sprintf("%%%s%%", *filter.UserName)})
	}

	return wherePred
}
