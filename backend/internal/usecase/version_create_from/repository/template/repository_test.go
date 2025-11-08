package template_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
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

		// project
		project := test_db.GenerateEntity(func(p *test_db.Project) {
			p.AuthorID = users[0].ID
		})
		projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

		// template
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.IsDefault = false
			t.ProjectID = &projectID
			t.AuthorID = &users[1].ID
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)

		want := domain.Template{
			AuthorID:        *template.AuthorID,
			ProjectAuthorID: project.AuthorID,
		}
		require.Equal(t, want, *got)
	})

	s.T().Run("IsDefault", func(t *testing.T) {
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.IsDefault = true
			t.ProjectID = nil
			t.AuthorID = nil
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
