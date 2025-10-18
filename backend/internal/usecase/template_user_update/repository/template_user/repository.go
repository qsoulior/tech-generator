package template_user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
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

func (r *Repository) GetByTemplateID(ctx context.Context, templateID int64) ([]domain.TemplateUser, error) {
	op := "template user - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"user_id",
			"role",
		).
		From("template_user").
		Where(sq.Eq{"template_id": templateID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var templateUsers []templateUser
	err = r.db.SelectContext(ctx, &templateUsers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	users := lo.Map(templateUsers, func(u templateUser, _ int) domain.TemplateUser { return u.toDomain() })
	return users, nil
}

func (r *Repository) Upsert(ctx context.Context, templateID int64, users []domain.TemplateUser) error {
	op := "template user - upsert"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("template_user").
		Columns("template_id", "user_id", "role").
		Suffix("ON CONFLICT (template_id, user_id) DO UPDATE SET role = EXCLUDED.role")

	for _, u := range users {
		builder = builder.Values(templateID, u.ID, u.Role)
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

func (r *Repository) Delete(ctx context.Context, templateID int64, userIDs []int64) error {
	op := "template user - delete"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("template_user").
		Where(sq.Eq{"template_id": templateID, "user_id": userIDs})

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
