package project_user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
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

func (r *Repository) GetByProjectID(ctx context.Context, projectID int64) ([]domain.ProjectUser, error) {
	op := "project user - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"user_id",
			"role",
		).
		From("project_user").
		Where(sq.Eq{"project_id": projectID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []projectUser
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	users := lo.Map(dtos, func(u projectUser, _ int) domain.ProjectUser { return u.toDomain() })
	return users, nil
}

func (r *Repository) Upsert(ctx context.Context, projectID int64, users []domain.ProjectUser) error {
	op := "project user - upsert"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("project_user").
		Columns("project_id", "user_id", "role").
		Suffix("ON CONFLICT (project_id, user_id) DO UPDATE SET role = EXCLUDED.role")

	for _, u := range users {
		builder = builder.Values(projectID, u.ID, u.Role)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query %q: %w", op, err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, projectID int64, userIDs []int64) error {
	op := "project user - delete"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("project_user").
		Where(sq.Eq{"project_id": projectID, "user_id": userIDs})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query %q: %w", op, err)
	}

	return nil
}
