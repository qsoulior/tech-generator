package user_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_create/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_Create() {
	ctx := context.Background()

	repo := New(s.C().DB())

	want := test_db.GenerateEntity[test_db.User]()

	s.T().Run("ExistsByName", func(t *testing.T) {
		user := test_db.GenerateEntity[test_db.User]()

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		err = repo.Create(ctx, user.Name, want.Email, want.Password)
		require.ErrorIs(t, err, domain.ErrUserExists)
	})

	s.T().Run("ExistsByEmail", func(t *testing.T) {
		user := test_db.GenerateEntity[test_db.User]()

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		err = repo.Create(ctx, want.Name, user.Email, want.Password)
		require.ErrorIs(t, err, domain.ErrUserExists)
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
