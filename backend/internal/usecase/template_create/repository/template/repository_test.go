package template_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
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

	// users
	users := test_db.GenerateEntities[test_db.User](2)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// folder
	folder := test_db.GenerateEntity(func(f *test_db.Folder) {
		f.ParentID = nil
		f.AuthorID = userIDs[0]
		f.RootAuthorID = userIDs[1]
	})
	folderID, err := test_db.InsertEntityWithID[int64](s.C(), "folder", folder)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "folder", folderID)) }()

	// template
	want := test_db.Template{
		Name:         gofakeit.UUID(),
		IsDefault:    false,
		FolderID:     &folderID,
		AuthorID:     &userIDs[1],
		RootAuthorID: &userIDs[1],
	}

	template := domain.Template{
		Name:         want.Name,
		IsDefault:    false,
		FolderID:     folderID,
		AuthorID:     userIDs[1],
		RootAuthorID: userIDs[1],
	}

	err = repo.Create(ctx, template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByColumn(s.C(), "template", "name", template.Name)) }()

	templates, err := test_db.SelectEntitiesByColumn[test_db.Template](s.C(), "template", "name", []string{template.Name})
	require.NoError(s.T(), err)
	require.Len(s.T(), templates, 1)

	got := templates[0]
	want.ID = got.ID
	want.CreatedAt = got.CreatedAt
	require.Equal(s.T(), want, got)
}
