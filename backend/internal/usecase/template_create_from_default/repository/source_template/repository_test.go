package source_template_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
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

	s.T().Run("Default", func(t *testing.T) {
		template := test_db.GenerateEntity(func(p *test_db.Template) {
			p.IsDefault = true
			p.ProjectID = nil
			p.AuthorID = nil
			p.LastVersionID = nil
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, templateID, got.ID)
		require.True(t, got.IsDefault)
		require.Nil(t, got.LastVersionID)
	})

	s.T().Run("NotDefault", func(t *testing.T) {
		user := test_db.GenerateEntity[test_db.User]()
		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		project := test_db.GenerateEntity(func(p *test_db.Project) { p.AuthorID = userID })
		projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

		template := test_db.GenerateEntity(func(p *test_db.Template) {
			p.IsDefault = false
			p.AuthorID = &userID
			p.ProjectID = &projectID
			p.LastVersionID = nil
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, templateID, got.ID)
		require.False(t, got.IsDefault)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}
