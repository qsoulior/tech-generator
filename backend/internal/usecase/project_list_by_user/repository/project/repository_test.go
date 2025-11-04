package project_repository

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_ListByUserID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	users := test_db.GenerateEntities[test_db.User](2)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// projects
	ownedProjects := test_db.GenerateEntities(2, func(p *test_db.Project, i int) { p.AuthorID = userIDs[0] })
	otherProjects := test_db.GenerateEntities(3, func(p *test_db.Project, i int) { p.AuthorID = userIDs[1] })
	projects := slices.Concat(ownedProjects, otherProjects)
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	// project users
	projectUsers := []test_db.ProjectUser{
		test_db.GenerateEntity(func(pu *test_db.ProjectUser) {
			pu.ProjectID = otherProjects[0].ID
			pu.UserID = userIDs[0]
		}),
		test_db.GenerateEntity(func(pu *test_db.ProjectUser) {
			pu.ProjectID = otherProjects[0].ID
			pu.UserID = userIDs[1]
		}),
	}
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "project_user", projectUsers, "user_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "project_user", "user_id", userIDs))
	}()

	wantProjects := slices.Concat(
		lo.Map(projects[:2], func(p test_db.Project, _ int) domain.Project {
			return domain.Project{
				ID:         p.ID,
				Name:       p.Name,
				AuthorName: users[0].Name,
			}
		}),
		lo.Map(projects[2:3], func(p test_db.Project, _ int) domain.Project {
			return domain.Project{
				ID:         p.ID,
				Name:       p.Name,
				AuthorName: users[1].Name,
			}
		}),
	)
	slices.SortFunc(wantProjects, func(a, b domain.Project) int { return int(b.ID - a.ID) })

	in := domain.ProjectListByUserIn{
		Page:    1,
		Size:    int64(len(projects)),
		Filter:  domain.ProjectListByUserFilter{UserID: userIDs[0]},
		Sorting: nil,
	}

	gotProjects, err := repo.ListByUserID(ctx, in)
	require.NoError(s.T(), err)
	require.Equal(s.T(), wantProjects, gotProjects)

	gotTotal, err := repo.GetTotalByUserID(ctx, in)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), len(gotProjects), gotTotal)
}

func (s *repositorySuite) TestRepository_ListByUserID_Filter() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	const N = 5

	// projects
	projects := test_db.GenerateEntities(N, func(p *test_db.Project, i int) {
		p.AuthorID = userID
	})
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	randomProject := lo.Sample(projects)

	tests := []struct {
		name       string
		filter     domain.ProjectListByUserFilter
		filterPred func(p domain.Project) bool
	}{
		{
			name: "ProjectName",
			filter: domain.ProjectListByUserFilter{
				UserID:      userID,
				ProjectName: lo.ToPtr(lo.Substring(randomProject.Name, 1, 3)),
			},
			filterPred: func(p domain.Project) bool {
				return p.AuthorName == user.Name && p.Name == randomProject.Name
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.ProjectListByUserIn{
				Page:    1,
				Size:    int64(len(projects)),
				Filter:  tt.filter,
				Sorting: nil,
			}

			gotProjects, err := repo.ListByUserID(ctx, in)
			require.NoError(s.T(), err)

			isFiltered := lo.EveryBy(gotProjects, tt.filterPred)
			require.True(t, isFiltered)

			gotTotal, err := repo.GetTotalByUserID(ctx, in)
			require.NoError(s.T(), err)
			require.EqualValues(t, len(gotProjects), gotTotal)
		})
	}
}

func (s *repositorySuite) TestRepository_ListByUserID_Sorting() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// projects
	projects := test_db.GenerateEntities(5, func(p *test_db.Project, i int) {
		p.AuthorID = userID
	})
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	tests := []struct {
		name       string
		sorting    *sorting_domain.Sorting
		comparator func(a, b domain.Project) int
	}{
		{
			name:    "Default",
			sorting: nil,
			comparator: func(a, b domain.Project) int {
				return int(b.ID - a.ID)
			},
		},
		{
			name: "Invalid",
			sorting: &sorting_domain.Sorting{
				Attribute: "invalid",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Project) int {
				return int(b.ID - a.ID)
			},
		},
		{
			name: "ProjectName_ASC",
			sorting: &sorting_domain.Sorting{
				Attribute: "project_name",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Project) int {
				return strings.Compare(a.Name, b.Name)
			},
		},
		{
			name: "ProjectName_DESC",
			sorting: &sorting_domain.Sorting{
				Attribute: "project_name",
				Direction: "DESC",
			},
			comparator: func(a, b domain.Project) int {
				return strings.Compare(b.Name, a.Name)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.ProjectListByUserIn{
				Page:    1,
				Size:    int64(len(projects)),
				Filter:  domain.ProjectListByUserFilter{UserID: userID},
				Sorting: tt.sorting,
			}

			got, err := repo.ListByUserID(ctx, in)
			require.NoError(s.T(), err)

			isSorted := slices.IsSortedFunc(got, tt.comparator)
			require.True(t, isSorted)
		})
	}
}

func (s *repositorySuite) TestRepository_ListByUserID_Pagination() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// projects
	projects := test_db.GenerateEntities(5, func(p *test_db.Project, i int) {
		p.AuthorID = userID
	})
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	wantProjects := lo.Map(projects, func(p test_db.Project, _ int) domain.Project {
		return domain.Project{
			ID:         p.ID,
			Name:       p.Name,
			AuthorName: user.Name,
		}
	})
	slices.SortFunc(wantProjects, func(a, b domain.Project) int { return int(b.ID - a.ID) })

	tests := []struct {
		name string
		page int64
		size int64
		want []domain.Project
	}{
		{
			name: "None",
			page: 2,
			size: int64(len(projects)),
			want: []domain.Project{},
		},
		{
			name: "One",
			page: 1,
			size: 1,
			want: wantProjects[:1],
		},
		{
			name: "Many",
			page: 2,
			size: 2,
			want: wantProjects[2:4],
		},
		{
			name: "All",
			page: 1,
			size: int64(len(projects)),
			want: wantProjects,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.ProjectListByUserIn{
				Page:    tt.page,
				Size:    tt.size,
				Filter:  domain.ProjectListByUserFilter{UserID: userID},
				Sorting: nil,
			}

			got, err := repo.ListByUserID(ctx, in)
			require.NoError(s.T(), err)
			require.Len(t, got, len(tt.want))
			require.Equal(t, tt.want, got)
		})
	}
}
