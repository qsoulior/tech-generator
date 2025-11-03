package project_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
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
		// user
		user := test_db.GenerateEntity[test_db.User]()
		userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

		// project
		project := test_db.GenerateEntity(func(p *test_db.Project) {
			p.AuthorID = userID
		})
		projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

		got, err := repo.GetByID(ctx, projectID)
		require.NoError(t, err)

		want := domain.Project{AuthorID: project.AuthorID}
		require.Equal(t, want, *got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}
