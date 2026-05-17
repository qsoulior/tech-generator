package user_repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

const pgUniqueViolation = "23505"

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, name, email string, password []byte) error {
	op := "usr - create"
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("usr").
		Columns("name", "email", "password").
		Values(name, email, password)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.ErrUserExists
		}
		return fmt.Errorf("exec query %q: %w", op, err)
	}

	return nil
}
