package project_user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByProjectID(ctx context.Context, projectID int64) ([]domain.ProjectUser, error) {
	op := "project user - get by project id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"u.id",
			"u.name",
			"u.email",
			"pu.role",
		).
		From("project_user pu").
		Join("usr u ON pu.user_id = u.id").
		Where(sq.Eq{"pu.project_id": projectID}).
		OrderBy("u.name ASC")

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
