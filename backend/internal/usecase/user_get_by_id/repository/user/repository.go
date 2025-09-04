package user_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	op := "usr - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"name",
			"email",
			"created_at",
		).
		From("usr").
		Where(sq.Eq{"id": id})

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
