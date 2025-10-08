package folder_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
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

		// folder users
		folderUsers := test_db.GenerateEntities(2, func(entity *test_db.FolderUser, i int) {
			entity.FolderID = folderID
			entity.UserID = userIDs[2:][i]
		})
		_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "folder_user", folderUsers, "folder_id")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, test_db.DeleteEntitiesByColumn(s.C(), "folder_user", "folder_id", []int64{folderID}))
		}()

		got, err := repo.GetByID(ctx, folderID)
		require.NoError(t, err)

		want := domain.Folder{
			AuthorID:     folder.AuthorID,
			RootAuthorID: folder.RootAuthorID,
			Users: []domain.FolderUser{
				{ID: folderUsers[0].UserID, Role: user_domain.Role(folderUsers[0].Role)},
				{ID: folderUsers[1].UserID, Role: user_domain.Role(folderUsers[1].Role)},
			},
		}
		require.Equal(t, want, *got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func (s *repositorySuite) TestRepository_Create() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	users := test_db.GenerateEntities[test_db.User](2)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// parent folder
	parent := test_db.GenerateEntity(func(f *test_db.Folder) {
		f.ParentID = nil
		f.AuthorID = userIDs[0]
		f.RootAuthorID = userIDs[1]
	})
	parentID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", parent)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "folder", parentID)) }()

	want := test_db.Folder{
		ParentID:     &parentID,
		Name:         gofakeit.UUID(),
		AuthorID:     userIDs[1],
		RootAuthorID: userIDs[1],
	}

	folder := domain.FolderToCreate{
		ParentID:     &parentID,
		Name:         want.Name,
		AuthorID:     userIDs[1],
		RootAuthorID: userIDs[1],
	}
	err = repo.Create(ctx, folder)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByColumn(s.C(), "folder", "author_id", userIDs[1])) }()

	folders, err := test_db.SelectEntitiesByColumn[test_db.Folder](s.C(), "folder", "author_id", []int64{userIDs[1]})
	require.NoError(s.T(), err)
	require.Len(s.T(), folders, 1)

	got := folders[0]
	want.ID = got.ID
	require.Equal(s.T(), want, got)
}
