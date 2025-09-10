package user_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByName(ctx context.Context, name string) (*domain.User, error) {
	op := "usr - get by name"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"password",
		).
		From("usr").
		Where(sq.Eq{"name": name})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var user user
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return user.toDomain(), nil
}
