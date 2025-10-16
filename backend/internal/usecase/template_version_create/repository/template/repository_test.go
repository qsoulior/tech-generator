package template_repository

import (
	"context"
	"testing"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByID() {
	ctx := context.Background()

	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

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

		// template
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.IsDefault = false
			t.FolderID = &folderID
			t.AuthorID = &users[0].ID
			t.RootAuthorID = &users[1].ID
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		// template users
		templateUsers := test_db.GenerateEntities(2, func(u *test_db.TemplateUser, i int) {
			u.TemplateID = templateID
			u.UserID = userIDs[2:][i]
		})
		_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "template_user", templateUsers, "template_id")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, test_db.DeleteEntitiesByColumn(s.C(), "template_user", "template_id", []int64{templateID}))
		}()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)

		want := domain.Template{
			AuthorID:     *template.AuthorID,
			RootAuthorID: *template.RootAuthorID,
			Users: []domain.TemplateUser{
				{ID: templateUsers[0].UserID, Role: user_domain.Role(templateUsers[0].Role)},
				{ID: templateUsers[1].UserID, Role: user_domain.Role(templateUsers[1].Role)},
			},
		}
		require.Equal(t, want, *got)
	})

	s.T().Run("IsDefault", func(t *testing.T) {
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.IsDefault = true
			t.FolderID = nil
			t.AuthorID = nil
			t.RootAuthorID = nil
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func (s *repositorySuite) TestRepository_UpdateByID() {
	ctx := context.Background()
	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// folder
	folder := test_db.GenerateEntity(func(f *test_db.Folder) {
		f.ParentID = nil
		f.AuthorID = userID
		f.RootAuthorID = userID
	})
	folderID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", folder)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "folder", folderID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.FolderID = &folderID
		t.AuthorID = &userID
		t.RootAuthorID = &userID
		t.CreatedAt = gofakeit.Date().Truncate(1 * time.Second)
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template updating
	templateToUpdate := domain.TemplateToUpdate{
		ID:            templateID,
		LastVersionID: gofakeit.Int64(),
	}
	err = repo.UpdateByID(ctx, templateToUpdate)
	require.NoError(s.T(), err)

	templates, err := test_db.SelectEntitiesByID[test_db.Template](s.C(), "template", []int64{templateID})
	require.NoError(s.T(), err)
	require.Len(s.T(), templates, 1)

	got := templates[0]
	require.Equal(s.T(), templateToUpdate.LastVersionID, *got.LastVersionID)

	now := time.Now().UTC().Truncate(1 * time.Second)
	require.GreaterOrEqual(s.T(), *got.UpdatedAt, now)
}
