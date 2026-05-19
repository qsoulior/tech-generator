package template_user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByTemplateID(ctx context.Context, templateID int64) ([]domain.TemplateUser, error) {
	op := "template user - get by template id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"u.id",
			"u.name",
			"u.email",
			"tu.role",
		).
		From("template_user tu").
		Join("usr u ON tu.user_id = u.id").
		Where(sq.Eq{"tu.template_id": templateID}).
		OrderBy("u.name ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []templateUser
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	users := lo.Map(dtos, func(u templateUser, _ int) domain.TemplateUser { return u.toDomain() })
	return users, nil
}
