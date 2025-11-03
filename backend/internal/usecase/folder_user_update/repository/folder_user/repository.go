package folder_user_repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
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

func (r *Repository) GetByFolderID(ctx context.Context, folderID int64) ([]domain.FolderUser, error) {
	op := "folder user - get by id"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"user_id",
			"role",
		).
		From("folder_user").
		Where(sq.Eq{"folder_id": folderID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query %q: %w", op, err)
	}

	query = fmt.Sprintf("-- %s\n%s", op, query)

	var dtos []folderUser
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query %q: %w", op, err)
	}

	users := lo.Map(dtos, func(u folderUser, _ int) domain.FolderUser { return u.toDomain() })
	return users, nil
}

func (r *Repository) Upsert(ctx context.Context, folderID int64, users []domain.FolderUser) error {
	op := "folder user - upsert"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("folder_user").
		Columns("folder_id", "user_id", "role").
		Suffix("ON CONFLICT (folder_id, user_id) DO UPDATE SET role = EXCLUDED.role")

	for _, u := range users {
		builder = builder.Values(folderID, u.ID, u.Role)
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

func (r *Repository) Delete(ctx context.Context, folderID int64, userIDs []int64) error {
	op := "folder user - delete"

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("folder_user").
		Where(sq.Eq{"folder_id": folderID, "user_id": userIDs})

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
