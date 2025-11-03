package folder_user_repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByFolderID() {
	ctx := context.Background()

	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// users
	users := test_db.GenerateEntities[test_db.User](4)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// folder
	folder := test_db.GenerateEntity(func(f *test_db.Folder) {
		f.ParentID = nil
		f.AuthorID = users[0].ID
		f.RootAuthorID = users[1].ID
	})
	folderID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", folder)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "folder", folderID)) }()

	// folder users
	folderUsers := test_db.GenerateEntities(2, func(u *test_db.FolderUser, i int) {
		u.FolderID = folderID
		u.UserID = userIDs[2:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "folder_user", folderUsers, "folder_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "folder_user", "folder_id", []int64{folderID}))
	}()

	got, err := repo.GetByFolderID(ctx, folderID)
	require.NoError(s.T(), err)

	want := []domain.FolderUser{
		{ID: folderUsers[0].UserID, Role: user_domain.Role(folderUsers[0].Role)},
		{ID: folderUsers[1].UserID, Role: user_domain.Role(folderUsers[1].Role)},
	}
	require.Equal(s.T(), want, got)
}

func (s *repositorySuite) TestRepository_Upsert() {
	ctx := context.Background()
	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// users
	users := test_db.GenerateEntities[test_db.User](3)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// folder
	folder := test_db.GenerateEntity(func(f *test_db.Folder) {
		f.ParentID = nil
		f.AuthorID = users[0].ID
		f.RootAuthorID = users[0].ID
	})
	folderID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", folder)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "folder", folderID)) }()

	// existing folder user
	folderUserExisting := test_db.FolderUser{
		FolderID: folderID,
		UserID:   users[1].ID,
		Role:     string(user_domain.RoleRead),
	}
	_, err = test_db.InsertEntityWithColumn[int64](s.C(), "folder_user", folderUserExisting, "folder_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "folder_user", "folder_id", []int64{folderID}))
	}()

	folderUsers := []domain.FolderUser{
		{ID: users[1].ID, Role: user_domain.RoleWrite},
		{ID: users[2].ID, Role: user_domain.RoleRead},
	}
	err = repo.Upsert(ctx, folderID, folderUsers)
	require.NoError(s.T(), err)

	got, err := test_db.SelectEntitiesByColumn[test_db.FolderUser](s.C(), "folder_user", "folder_id", []int64{folderID})
	require.NoError(s.T(), err)

	want := []test_db.FolderUser{
		{
			FolderID: folderID,
			UserID:   users[1].ID,
			Role:     string(user_domain.RoleWrite),
		},
		{
			FolderID: folderID,
			UserID:   users[2].ID,
			Role:     string(user_domain.RoleRead),
		},
	}
	require.Equal(s.T(), want, got)
}

func (s *repositorySuite) TestRepository_Delete() {
	ctx := context.Background()
	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// users
	users := test_db.GenerateEntities[test_db.User](4)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// folder
	folder := test_db.GenerateEntity(func(f *test_db.Folder) {
		f.ParentID = nil
		f.AuthorID = users[0].ID
		f.RootAuthorID = users[0].ID
	})
	folderID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", folder)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "folder", folderID)) }()

	// folder users
	folderUsers := test_db.GenerateEntities(3, func(u *test_db.FolderUser, i int) {
		u.FolderID = folderID
		u.UserID = userIDs[1:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "folder_user", folderUsers, "folder_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "folder_user", "folder_id", []int64{folderID}))
	}()

	err = repo.Delete(ctx, folderID, userIDs[1:3])
	require.NoError(s.T(), err)

	got, err := test_db.SelectEntitiesByColumn[test_db.FolderUser](s.C(), "folder_user", "folder_id", []int64{folderID})
	require.NoError(s.T(), err)

	want := folderUsers[2:3]
	require.Equal(s.T(), want, got)
}
