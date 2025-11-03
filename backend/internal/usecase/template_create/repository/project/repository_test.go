package project_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
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

		// project users
		projectUsers := test_db.GenerateEntities(2, func(u *test_db.ProjectUser, i int) {
			u.ProjectID = projectID
			u.UserID = userIDs[2:][i]
		})
		_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "project_user", projectUsers, "project_id")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, test_db.DeleteEntitiesByColumn(s.C(), "project_user", "project_id", []int64{projectID}))
		}()

		got, err := repo.GetByID(ctx, projectID)
		require.NoError(t, err)

		want := domain.Project{
			AuthorID: project.AuthorID,
			Users: []domain.ProjectUser{
				{ID: projectUsers[0].UserID, Role: user_domain.Role(projectUsers[0].Role)},
				{ID: projectUsers[1].UserID, Role: user_domain.Role(projectUsers[1].Role)},
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
