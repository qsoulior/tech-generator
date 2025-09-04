package user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

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
		return fmt.Errorf("exec query %q: %w", op, err)
	}

	return nil
}

func (r *Repository) ExistsByNameOrEmail(ctx context.Context, name, email string) (bool, error) {
	op := "usr - exists by name or email"
	selectBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("1").
		From("usr").
		Where(sq.Or{sq.Eq{"name": name}, sq.Eq{"email": email}})

	existsBuilder := selectBuilder.
		Prefix("SELECT EXISTS (").
		Suffix(")")

	query, args, err := existsBuilder.ToSql()
	if err != nil {
		return false, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var exists bool
	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("exec query %q: %w", op, err)
	}

	return exists, nil
}
