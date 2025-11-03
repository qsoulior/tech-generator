package project_repository

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_ListByAuthorID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	users := test_db.GenerateEntities[test_db.User](2)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// projects
	projects := slices.Concat(
		// owned
		test_db.GenerateEntities(3, func(p *test_db.Project, i int) {
			p.AuthorID = userIDs[0]
		}),
		// other
		test_db.GenerateEntities(2, func(p *test_db.Project, i int) {
			p.AuthorID = userIDs[1]
		}),
	)
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	wantProjects := lo.Map(projects[:3], func(p test_db.Project, _ int) domain.Project {
		return domain.Project{
			ID:   p.ID,
			Name: p.Name,
		}
	})
	slices.SortFunc(wantProjects, func(a, b domain.Project) int { return strings.Compare(a.Name, b.Name) })

	tests := []struct {
		name string
		in   domain.ProjectListByUserIn
		want []domain.Project
	}{
		{
			name: "All",
			in: domain.ProjectListByUserIn{
				UserID: userIDs[0],
				Page:   1,
				Size:   uint64(len(projects)),
			},
			want: wantProjects,
		},
		{
			name: "Part",
			in: domain.ProjectListByUserIn{
				UserID: userIDs[0],
				Page:   2,
				Size:   2,
			},
			want: wantProjects[2:],
		},
		{
			name: "None",
			in: domain.ProjectListByUserIn{
				UserID: userIDs[0],
				Page:   3,
				Size:   2,
			},
			want: []domain.Project{},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := repo.ListByAuthorID(ctx, tt.in)
			require.NoError(s.T(), err)
			require.Equal(t, tt.want, got)
		})
	}
}

func (s *repositorySuite) TestRepository_ListByProjectUserID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	users := test_db.GenerateEntities[test_db.User](3)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// projects
	projects := test_db.GenerateEntities(5, func(p *test_db.Project, i int) {
		p.AuthorID = userIDs[2]
	})
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	// project users
	projectUsers := slices.Concat(
		// shared
		test_db.GenerateEntities(3, func(pu *test_db.ProjectUser, i int) {
			pu.ProjectID = projectIDs[i]
			pu.UserID = userIDs[0]
		}),
		// other
		test_db.GenerateEntities(2, func(pu *test_db.ProjectUser, i int) {
			pu.ProjectID = projectIDs[i]
			pu.UserID = userIDs[1]
		}),
	)
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "project_user", projectUsers, "user_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "project_user", "user_id", userIDs))
	}()

	wantProjects := lo.Map(projects[:3], func(p test_db.Project, _ int) domain.Project {
		return domain.Project{
			ID:   p.ID,
			Name: p.Name,
		}
	})
	slices.SortFunc(wantProjects, func(a, b domain.Project) int { return strings.Compare(a.Name, b.Name) })

	tests := []struct {
		name string
		in   domain.ProjectListByUserIn
		want []domain.Project
	}{
		{
			name: "All",
			in: domain.ProjectListByUserIn{
				UserID: userIDs[0],
				Page:   1,
				Size:   uint64(len(projects)),
			},
			want: wantProjects,
		},
		{
			name: "Part",
			in: domain.ProjectListByUserIn{
				UserID: userIDs[0],
				Page:   2,
				Size:   2,
			},
			want: wantProjects[2:],
		},
		{
			name: "None",
			in: domain.ProjectListByUserIn{
				UserID: userIDs[0],
				Page:   3,
				Size:   2,
			},
			want: []domain.Project{},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := repo.ListByProjectUserID(ctx, tt.in)
			require.NoError(s.T(), err)
			require.Equal(t, tt.want, got)
		})
	}
}
