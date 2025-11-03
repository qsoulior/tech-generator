package project_user_repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByProjectID() {
	ctx := context.Background()

	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// users
	users := test_db.GenerateEntities[test_db.User](4)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) {
		p.AuthorID = users[0].ID
	})
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// project users
	projectUsers := test_db.GenerateEntities(2, func(u *test_db.ProjectUser, i int) {
		u.ProjectID = projectID
		u.UserID = userIDs[2:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "project_user", projectUsers, "project_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "project_user", "project_id", []int64{projectID}))
	}()

	got, err := repo.GetByProjectID(ctx, projectID)
	require.NoError(s.T(), err)

	want := []domain.ProjectUser{
		{ID: projectUsers[0].UserID, Role: user_domain.Role(projectUsers[0].Role)},
		{ID: projectUsers[1].UserID, Role: user_domain.Role(projectUsers[1].Role)},
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

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) {
		p.AuthorID = users[0].ID
	})
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// existing project user
	projectUserExisting := test_db.ProjectUser{
		ProjectID: projectID,
		UserID:    users[1].ID,
		Role:      string(user_domain.RoleRead),
	}
	_, err = test_db.InsertEntityWithColumn[int64](s.C(), "project_user", projectUserExisting, "project_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "project_user", "project_id", []int64{projectID}))
	}()

	projectUsers := []domain.ProjectUser{
		{ID: users[1].ID, Role: user_domain.RoleWrite},
		{ID: users[2].ID, Role: user_domain.RoleRead},
	}
	err = repo.Upsert(ctx, projectID, projectUsers)
	require.NoError(s.T(), err)

	got, err := test_db.SelectEntitiesByColumn[test_db.ProjectUser](s.C(), "project_user", "project_id", []int64{projectID})
	require.NoError(s.T(), err)

	want := []test_db.ProjectUser{
		{
			ProjectID: projectID,
			UserID:    users[1].ID,
			Role:      string(user_domain.RoleWrite),
		},
		{
			ProjectID: projectID,
			UserID:    users[2].ID,
			Role:      string(user_domain.RoleRead),
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

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) {
		p.AuthorID = users[0].ID
	})
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// project users
	projectUsers := test_db.GenerateEntities(3, func(u *test_db.ProjectUser, i int) {
		u.ProjectID = projectID
		u.UserID = userIDs[1:][i]
	})
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "project_user", projectUsers, "project_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "project_user", "project_id", []int64{projectID}))
	}()

	err = repo.Delete(ctx, projectID, userIDs[1:3])
	require.NoError(s.T(), err)

	got, err := test_db.SelectEntitiesByColumn[test_db.ProjectUser](s.C(), "project_user", "project_id", []int64{projectID})
	require.NoError(s.T(), err)

	want := projectUsers[2:3]
	require.Equal(s.T(), want, got)
}
