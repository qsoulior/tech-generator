package template_repository

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
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
	projects := test_db.GenerateEntities(2, func(p *test_db.Project, i int) {
		p.AuthorID = userIDs[1]
	})
	projectIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "project", projects)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "project", projectIDs)) }()

	// templates
	templates := slices.Concat(
		// owned, in project
		test_db.GenerateEntities(3, func(p *test_db.Template, i int) {
			p.AuthorID = &userIDs[0]
			p.ProjectID = &projectIDs[0]
		}),
		// owned, not in project
		test_db.GenerateEntities(2, func(p *test_db.Template, i int) {
			p.AuthorID = &userIDs[0]
			p.ProjectID = &projectIDs[1]
		}),
		// other, in project
		test_db.GenerateEntities(2, func(p *test_db.Template, i int) {
			p.AuthorID = &userIDs[1]
			p.ProjectID = &projectIDs[0]
		}),
	)
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	// template users
	templateUsers := []test_db.TemplateUser{
		test_db.GenerateEntity(func(pu *test_db.TemplateUser) {
			pu.TemplateID = templates[5:][0].ID
			pu.UserID = userIDs[0]
		}),
		test_db.GenerateEntity(func(pu *test_db.TemplateUser) {
			pu.TemplateID = templates[5:][0].ID
			pu.UserID = userIDs[1]
		}),
	}
	_, err = test_db.InsertEntitiesWithColumn[int64](s.C(), "template_user", templateUsers, "user_id")
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "template_user", "user_id", userIDs))
	}()

	wantTemplates := slices.Concat(
		lo.Map(templates[:3], func(p test_db.Template, _ int) domain.Template {
			return domain.Template{
				ID:         p.ID,
				Name:       p.Name,
				AuthorName: users[0].Name,
				CreatedAt:  p.CreatedAt.Truncate(1 * time.Microsecond),
				UpdatedAt:  lo.ToPtr(p.UpdatedAt.Truncate(1 * time.Microsecond)),
			}
		}),
		lo.Map(templates[5:6], func(p test_db.Template, _ int) domain.Template {
			return domain.Template{
				ID:         p.ID,
				Name:       p.Name,
				AuthorName: users[1].Name,
				CreatedAt:  p.CreatedAt.Truncate(1 * time.Microsecond),
				UpdatedAt:  lo.ToPtr(p.UpdatedAt.Truncate(1 * time.Microsecond)),
			}
		}),
	)
	slices.SortFunc(wantTemplates, func(a, b domain.Template) int { return int(b.ID - a.ID) })

	in := domain.TemplateListByUserIn{
		Page:    1,
		Size:    int64(len(templates)),
		Filter:  domain.TemplateListByUserFilter{UserID: userIDs[0], ProjectID: projectIDs[0]},
		Sorting: nil,
	}

	gotTemplates, err := repo.ListByUserID(ctx, in)
	require.NoError(s.T(), err)
	require.Equal(s.T(), wantTemplates, gotTemplates)

	gotTotal, err := repo.GetTotalByUserID(ctx, in)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), len(gotTemplates), gotTotal)
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

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) { p.AuthorID = userID })
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// templates
	templates := test_db.GenerateEntities(N, func(p *test_db.Template, i int) {
		p.AuthorID = &userID
		p.ProjectID = &projectID
	})
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	randomTemplate := lo.Sample(templates)

	tests := []struct {
		name       string
		filter     domain.TemplateListByUserFilter
		filterPred func(t domain.Template) bool
	}{
		{
			name: "TemplateName",
			filter: domain.TemplateListByUserFilter{
				UserID:       userID,
				ProjectID:    projectID,
				TemplateName: lo.ToPtr(lo.Substring(randomTemplate.Name, 1, 3)),
			},
			filterPred: func(t domain.Template) bool {
				return t.AuthorName == user.Name && t.Name == randomTemplate.Name
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TemplateListByUserIn{
				Page:    1,
				Size:    int64(len(templates)),
				Filter:  tt.filter,
				Sorting: nil,
			}

			gotTemplates, err := repo.ListByUserID(ctx, in)
			require.NoError(s.T(), err)

			isFiltered := lo.EveryBy(gotTemplates, tt.filterPred)
			require.True(t, isFiltered)

			gotTotal, err := repo.GetTotalByUserID(ctx, in)
			require.NoError(s.T(), err)
			require.EqualValues(t, len(gotTemplates), gotTotal)
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

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) { p.AuthorID = userID })
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// templates
	templates := test_db.GenerateEntities(5, func(p *test_db.Template, i int) {
		p.AuthorID = &userID
		p.ProjectID = &projectID
	})
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	tests := []struct {
		name       string
		sorting    *sorting_domain.Sorting
		comparator func(a, b domain.Template) int
	}{
		{
			name:    "Default",
			sorting: nil,
			comparator: func(a, b domain.Template) int {
				return int(b.ID - a.ID)
			},
		},
		{
			name: "Invalid",
			sorting: &sorting_domain.Sorting{
				Attribute: "invalid",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Template) int {
				return int(b.ID - a.ID)
			},
		},
		{
			name: "TemplateName_ASC",
			sorting: &sorting_domain.Sorting{
				Attribute: "template_name",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Template) int {
				return strings.Compare(a.Name, b.Name)
			},
		},
		{
			name: "TemplateName_DESC",
			sorting: &sorting_domain.Sorting{
				Attribute: "template_name",
				Direction: "DESC",
			},
			comparator: func(a, b domain.Template) int {
				return strings.Compare(b.Name, a.Name)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TemplateListByUserIn{
				Page:    1,
				Size:    int64(len(templates)),
				Filter:  domain.TemplateListByUserFilter{UserID: userID, ProjectID: projectID},
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

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) { p.AuthorID = userID })
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// templates
	templates := test_db.GenerateEntities(5, func(p *test_db.Template, i int) {
		p.AuthorID = &userID
		p.ProjectID = &projectID
	})
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	wantTemplates := lo.Map(templates, func(t test_db.Template, _ int) domain.Template {
		return domain.Template{
			ID:         t.ID,
			Name:       t.Name,
			AuthorName: user.Name,
			CreatedAt:  t.CreatedAt.Truncate(1 * time.Microsecond),
			UpdatedAt:  lo.ToPtr(t.UpdatedAt.Truncate(1 * time.Microsecond)),
		}
	})
	slices.SortFunc(wantTemplates, func(a, b domain.Template) int { return int(b.ID - a.ID) })

	tests := []struct {
		name string
		page int64
		size int64
		want []domain.Template
	}{
		{
			name: "None",
			page: 2,
			size: int64(len(templates)),
			want: []domain.Template{},
		},
		{
			name: "One",
			page: 1,
			size: 1,
			want: wantTemplates[:1],
		},
		{
			name: "Many",
			page: 2,
			size: 2,
			want: wantTemplates[2:4],
		},
		{
			name: "All",
			page: 1,
			size: int64(len(templates)),
			want: wantTemplates,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TemplateListByUserIn{
				Page: tt.page,
				Size: tt.size,
				Filter: domain.TemplateListByUserFilter{
					UserID:    userID,
					ProjectID: projectID,
				},
				Sorting: nil,
			}

			got, err := repo.ListByUserID(ctx, in)
			require.NoError(s.T(), err)
			require.Len(t, got, len(tt.want))
			require.Equal(t, tt.want, got)
		})
	}
}
