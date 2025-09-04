package user_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgconn"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const codeUniqueViolation = "23505"

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_Create() {
	ctx := context.Background()

	repo := New(s.C().DB())

	want, err := test_db.GenerateEntity[test_db.User]()
	require.NoError(s.T(), err)

	s.T().Run("ExistsByName", func(t *testing.T) {
		user, err := test_db.GenerateEntity[test_db.User]()
		require.NoError(t, err)

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		err = repo.Create(ctx, user.Name, want.Email, want.Password)

		var pgErr *pgconn.PgError
		require.ErrorAs(t, err, &pgErr)
		require.Equal(t, codeUniqueViolation, pgErr.Code)
	})

	s.T().Run("ExistsByEmail", func(t *testing.T) {
		user, err := test_db.GenerateEntity[test_db.User]()
		require.NoError(t, err)

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		err = repo.Create(ctx, want.Name, user.Email, want.Password)

		var pgErr *pgconn.PgError
		require.ErrorAs(t, err, &pgErr)
		require.Equal(t, codeUniqueViolation, pgErr.Code)
	})

	s.T().Run("DoesNotExist", func(t *testing.T) {
		err := repo.Create(ctx, want.Name, want.Email, want.Password)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByColumn(s.C(), "usr", "name", want.Name)) }()

		got, err := test_db.SelectEntitiesByColumn[test_db.User](s.C(), "usr", "name", []string{want.Name})
		require.NoError(t, err)
		require.Len(t, got, 1)

		require.Equal(t, want.Name, got[0].Name)
		require.Equal(t, want.Email, got[0].Email)
		require.Equal(t, want.Password, got[0].Password)
	})

}

func (s *repositorySuite) TestRepository_ExistsByNameOrEmail() {
	ctx := context.Background()
	repo := New(s.C().DB())

	s.T().Run("ExistsByName", func(t *testing.T) {
		user, err := test_db.GenerateEntity[test_db.User]()
		require.NoError(t, err)

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(s.T(), err)
		defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		got, err := repo.ExistsByNameOrEmail(ctx, user.Name, gofakeit.Email())
		require.NoError(t, err)
		require.True(t, got)
	})

	s.T().Run("ExistsByEmail", func(t *testing.T) {
		user, err := test_db.GenerateEntity[test_db.User]()
		require.NoError(t, err)

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(s.T(), err)
		defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		got, err := repo.ExistsByNameOrEmail(ctx, gofakeit.Username(), user.Email)
		require.NoError(t, err)
		require.True(t, got)
	})

	s.T().Run("DoesNotExist", func(t *testing.T) {
		got, err := repo.ExistsByNameOrEmail(ctx, gofakeit.Username(), gofakeit.Email())
		require.NoError(t, err)
		require.False(t, got)
	})
}
