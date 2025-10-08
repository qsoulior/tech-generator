package folder_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Folder, error) {
	op := "folder - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"f.author_id",
			"f.root_author_id",
			"fu.user_id",
			"fu.role",
		).
		From("folder f").
		LeftJoin("folder_user fu ON f.id = fu.folder_id").
		Where(sq.Eq{"f.id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var folders folders
	err = r.db.SelectContext(ctx, &folders, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	return folders.toDomain(), nil
}

func (r *Repository) Create(ctx context.Context, name string, authorID int64, rootAuthorID int64, parentID *int64) error {
	op := "folder - create"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("folder").
		Columns("parent_id", "name", "author_id", "root_author_id").
		Values(parentID, name, authorID, rootAuthorID)

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
