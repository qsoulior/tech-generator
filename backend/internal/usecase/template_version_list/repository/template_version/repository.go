package template_version_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_list/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ListByTemplateID(ctx context.Context, templateID int64) ([]domain.TemplateVersion, error) {
	op := "template version - list by template id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"v.id",
			"v.number",
			"u.name as author_name",
			"v.created_at",
		).
		From("template_version v").
		Join("usr u ON v.author_id = u.id").
		Where(sq.Eq{"v.template_id": templateID}).
		OrderBy("v.id DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []templateVersion
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	versions := lo.Map(dtos, func(dto templateVersion, _ int) domain.TemplateVersion { return dto.toDomain() })
	return versions, nil
}
