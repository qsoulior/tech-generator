package template_user_repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByTemplateID() {
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

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.FolderID = &folderID
		t.AuthorID = &users[0].ID
		t.RootAuthorID = &users[1].ID
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template users
	templateUsers := test_db.GenerateEntities(2, func(u *test_db.TemplateUser, i int) {
		u.TemplateID = templateID
		u.UserID = userIDs[2:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "template_user", templateUsers, "template_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "template_user", "template_id", []int64{templateID}))
	}()

	got, err := repo.GetByTemplateID(ctx, templateID)
	require.NoError(s.T(), err)

	want := []domain.TemplateUser{
		{ID: templateUsers[0].UserID, Role: user_domain.Role(templateUsers[0].Role)},
		{ID: templateUsers[1].UserID, Role: user_domain.Role(templateUsers[1].Role)},
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

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.FolderID = &folderID
		t.AuthorID = &users[0].ID
		t.RootAuthorID = &users[0].ID
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// existing template user
	templateUserExisting := test_db.TemplateUser{
		TemplateID: templateID,
		UserID:     users[1].ID,
		Role:       string(user_domain.RoleRead),
	}
	_, err = test_db.InsertEntityWithColumn[int64](s.C(), "template_user", templateUserExisting, "template_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "template_user", "template_id", []int64{templateID}))
	}()

	templateUsers := []domain.TemplateUser{
		{ID: users[1].ID, Role: user_domain.RoleWrite},
		{ID: users[2].ID, Role: user_domain.RoleRead},
	}
	err = repo.Upsert(ctx, templateID, templateUsers)
	require.NoError(s.T(), err)

	got, err := test_db.SelectEntitiesByColumn[test_db.TemplateUser](s.C(), "template_user", "template_id", []int64{templateID})
	require.NoError(s.T(), err)

	want := []test_db.TemplateUser{
		{
			TemplateID: templateID,
			UserID:     users[1].ID,
			Role:       string(user_domain.RoleWrite),
		},
		{
			TemplateID: templateID,
			UserID:     users[2].ID,
			Role:       string(user_domain.RoleRead),
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

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.FolderID = &folderID
		t.AuthorID = &users[0].ID
		t.RootAuthorID = &users[0].ID
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template users
	templateUsers := test_db.GenerateEntities(3, func(u *test_db.TemplateUser, i int) {
		u.TemplateID = templateID
		u.UserID = userIDs[1:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "template_user", templateUsers, "template_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "template_user", "template_id", []int64{templateID}))
	}()

	err = repo.Delete(ctx, templateID, userIDs[1:3])
	require.NoError(s.T(), err)

	got, err := test_db.SelectEntitiesByColumn[test_db.TemplateUser](s.C(), "template_user", "template_id", []int64{templateID})
	require.NoError(s.T(), err)

	want := templateUsers[2:3]
	require.Equal(s.T(), want, got)
}
