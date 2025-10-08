package user_repository

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByName() {
	ctx := context.Background()
	repo := New(s.C().DB())

	s.T().Run("Exists", func(t *testing.T) {
		user := test_db.GenerateEntity[test_db.User]()
		user.CreatedAt = user.CreatedAt.Truncate(1 * time.Second)

		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(s.T(), err)
		defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		got, err := repo.GetByName(ctx, user.Name)
		require.NoError(t, err)
		require.NotNil(t, got)

		want := testModelToDomain(user)
		require.Equal(t, want, *got)
	})

	s.T().Run("DoesNotExist", func(t *testing.T) {
		got, err := repo.GetByName(ctx, gofakeit.Username())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func testModelToDomain(user test_db.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Password: user.Password,
	}
}
