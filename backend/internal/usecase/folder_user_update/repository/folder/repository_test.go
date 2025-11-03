package folder_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByID() {
	ctx := context.Background()

	repo := New(s.C().DB())

	s.T().Run("Exists", func(t *testing.T) {
		// users
		users := test_db.GenerateEntities[test_db.User](4)
		userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

		// folder
		folder := test_db.GenerateEntity(func(f *test_db.Folder) {
			f.ParentID = nil
			f.AuthorID = users[0].ID
			f.RootAuthorID = users[1].ID
		})
		folderID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", folder)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "folder", folderID)) }()

		got, err := repo.GetByID(ctx, folderID)
		require.NoError(t, err)

		want := domain.Folder{
			AuthorID:     folder.AuthorID,
			RootAuthorID: folder.RootAuthorID,
		}
		require.Equal(t, want, *got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}
