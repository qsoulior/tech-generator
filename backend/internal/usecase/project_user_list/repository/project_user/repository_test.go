package project_user_repository

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByProjectID() {
	ctx := context.Background()

	repo := New(s.C().DB())

	users := test_db.GenerateEntities[test_db.User](3)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	project := test_db.GenerateEntity(func(p *test_db.Project) {
		p.AuthorID = users[0].ID
	})
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	projectUsers := test_db.GenerateEntities(2, func(u *test_db.ProjectUser, i int) {
		u.ProjectID = projectID
		u.UserID = userIDs[1:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "project_user", projectUsers, "project_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "project_user", "project_id", []int64{projectID}))
	}()

	got, err := repo.GetByProjectID(ctx, projectID)
	require.NoError(s.T(), err)

	want := []domain.ProjectUser{
		{ID: users[1].ID, Name: users[1].Name, Email: users[1].Email, Role: user_domain.Role(projectUsers[0].Role)},
		{ID: users[2].ID, Name: users[2].Name, Email: users[2].Email, Role: user_domain.Role(projectUsers[1].Role)},
	}
	sort.Slice(want, func(i, j int) bool { return want[i].Name < want[j].Name })
	require.Equal(s.T(), want, got)
}
