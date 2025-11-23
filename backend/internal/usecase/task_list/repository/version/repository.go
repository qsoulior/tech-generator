package version_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

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

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Version, error) {
	op := "version - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"p.author_id as project_author_id",
			"t.author_id as template_author_id",
			"tu.user_id as template_user_id",
			"tu.role as template_user_role",
		).
		From("template_version v").
		Join("template t ON v.template_id = t.id").
		Join("project p ON t.project_id = p.id").
		LeftJoin("template_user tu ON t.id = tu.template_id").
		Where(sq.Eq{"v.id": id, "t.is_default": false})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos versions
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return dtos.toDomain(), nil
}
