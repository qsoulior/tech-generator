package user_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/user_list/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_List() {
	ctx := context.Background()
	repo := New(s.C().DB())

	users := []test_db.User{
		test_db.GenerateEntity(func(u *test_db.User) { u.Name = "user_alpha" }),
		test_db.GenerateEntity(func(u *test_db.User) { u.Name = "user_beta" }),
		test_db.GenerateEntity(func(u *test_db.User) { u.Name = "other_gamma" }),
	}
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	s.T().Run("ExcludeSelfAndFilterByName", func(t *testing.T) {
		userName := "user_"
		got, err := repo.List(ctx, domain.UserListIn{
			Page: 1,
			Size: 10,
			Filter: domain.UserListFilter{
				ExcludeUserID: userIDs[0],
				UserName:      &userName,
			},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "user_beta", got[0].Name)
	})

	s.T().Run("Pagination", func(t *testing.T) {
		page1, err := repo.List(ctx, domain.UserListIn{
			Page:   1,
			Size:   1,
			Filter: domain.UserListFilter{ExcludeUserID: userIDs[0]},
		})
		require.NoError(t, err)
		require.Len(t, page1, 1)

		page2, err := repo.List(ctx, domain.UserListIn{
			Page:   2,
			Size:   1,
			Filter: domain.UserListFilter{ExcludeUserID: userIDs[0]},
		})
		require.NoError(t, err)
		require.Len(t, page2, 1)
		require.NotEqual(t, page1[0].ID, page2[0].ID)
	})
}

func (s *repositorySuite) TestRepository_GetTotal() {
	ctx := context.Background()
	repo := New(s.C().DB())

	users := []test_db.User{
		test_db.GenerateEntity(func(u *test_db.User) { u.Name = "user_alpha" }),
		test_db.GenerateEntity(func(u *test_db.User) { u.Name = "user_beta" }),
		test_db.GenerateEntity(func(u *test_db.User) { u.Name = "other_gamma" }),
	}
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	s.T().Run("OnlyExcludeSelf", func(t *testing.T) {
		got, err := repo.GetTotal(ctx, domain.UserListIn{
			Filter: domain.UserListFilter{ExcludeUserID: userIDs[0]},
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, got, int64(2))
	})

	s.T().Run("FilterByName", func(t *testing.T) {
		userName := "user_"
		got, err := repo.GetTotal(ctx, domain.UserListIn{
			Filter: domain.UserListFilter{
				ExcludeUserID: userIDs[0],
				UserName:      &userName,
			},
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), got)
	})
}
