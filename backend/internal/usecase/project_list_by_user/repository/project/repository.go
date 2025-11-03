package project_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

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

func (r *Repository) ListByAuthorID(ctx context.Context, in domain.ProjectListByUserIn) ([]domain.Project, error) {
	op := "project - list by author id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"name",
		).
		From("project").
		Where(sq.Eq{"author_id": in.UserID}).
		OrderBy("name ASC").
		Offset((in.Page - 1) * in.Size).
		Limit(in.Size)

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

func (r *Repository) ListByProjectUserID(ctx context.Context, in domain.ProjectListByUserIn) ([]domain.Project, error) {
	op := "project - list by author id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"p.id",
			"p.name",
		).
		From("project p").
		Join("project_user pu ON p.id = pu.project_id AND pu.user_id = ?", in.UserID).
		OrderBy("name ASC").
		Offset((in.Page - 1) * in.Size).
		Limit(in.Size)

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
