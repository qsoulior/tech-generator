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
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_ListDefault() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// default templates
	defaults := test_db.GenerateEntities(3, func(p *test_db.Template, i int) {
		p.IsDefault = true
		p.ProjectID = nil
		p.AuthorID = nil
	})
	defaultIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", defaults)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", defaultIDs)) }()

	// non-default templates: must be ignored, but need user+project
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	project := test_db.GenerateEntity(func(p *test_db.Project) { p.AuthorID = userID })
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	nonDefaults := test_db.GenerateEntities(2, func(p *test_db.Template, i int) {
		p.IsDefault = false
		p.AuthorID = &userID
		p.ProjectID = &projectID
	})
	nonDefaultIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", nonDefaults)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", nonDefaultIDs)) }()

	wantTemplates := lo.Map(defaults, func(p test_db.Template, _ int) domain.Template {
		return domain.Template{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt.Truncate(1 * time.Microsecond),
			UpdatedAt: lo.ToPtr(p.UpdatedAt.Truncate(1 * time.Microsecond)),
		}
	})
	slices.SortFunc(wantTemplates, func(a, b domain.Template) int { return int(b.ID - a.ID) })

	in := domain.TemplateListDefaultIn{
		Page:    1,
		Size:    int64(len(defaults) + len(nonDefaults)),
		Filter:  domain.TemplateListDefaultFilter{},
		Sorting: nil,
	}

	gotTemplates, err := repo.ListDefault(ctx, in)
	require.NoError(s.T(), err)
	require.Equal(s.T(), wantTemplates, gotTemplates)

	gotTotal, err := repo.GetTotalDefault(ctx, in)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), len(wantTemplates), gotTotal)
}

func (s *repositorySuite) TestRepository_ListDefault_Filter() {
	ctx := context.Background()
	repo := New(s.C().DB())

	const N = 5

	templates := test_db.GenerateEntities(N, func(p *test_db.Template, i int) {
		p.IsDefault = true
		p.ProjectID = nil
		p.AuthorID = nil
	})
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	randomTemplate := lo.Sample(templates)

	in := domain.TemplateListDefaultIn{
		Page: 1,
		Size: int64(N),
		Filter: domain.TemplateListDefaultFilter{
			TemplateName: lo.ToPtr(lo.Substring(randomTemplate.Name, 1, 3)),
		},
		Sorting: nil,
	}

	gotTemplates, err := repo.ListDefault(ctx, in)
	require.NoError(s.T(), err)

	isFiltered := lo.EveryBy(gotTemplates, func(t domain.Template) bool { return t.Name == randomTemplate.Name })
	require.True(s.T(), isFiltered)

	gotTotal, err := repo.GetTotalDefault(ctx, in)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), len(gotTemplates), gotTotal)
}

func (s *repositorySuite) TestRepository_ListDefault_Sorting() {
	ctx := context.Background()
	repo := New(s.C().DB())

	templates := test_db.GenerateEntities(5, func(p *test_db.Template, i int) {
		p.IsDefault = true
		p.ProjectID = nil
		p.AuthorID = nil
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
			name: "Name_ASC",
			sorting: &sorting_domain.Sorting{
				Attribute: "name",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Template) int {
				return strings.Compare(a.Name, b.Name)
			},
		},
		{
			name: "Name_DESC",
			sorting: &sorting_domain.Sorting{
				Attribute: "name",
				Direction: "DESC",
			},
			comparator: func(a, b domain.Template) int {
				return strings.Compare(b.Name, a.Name)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TemplateListDefaultIn{
				Page:    1,
				Size:    int64(len(templates)),
				Filter:  domain.TemplateListDefaultFilter{},
				Sorting: tt.sorting,
			}

			got, err := repo.ListDefault(ctx, in)
			require.NoError(s.T(), err)

			isSorted := slices.IsSortedFunc(got, tt.comparator)
			require.True(t, isSorted)
		})
	}
}

func (s *repositorySuite) TestRepository_ListDefault_Pagination() {
	ctx := context.Background()
	repo := New(s.C().DB())

	templates := test_db.GenerateEntities(5, func(p *test_db.Template, i int) {
		p.IsDefault = true
		p.ProjectID = nil
		p.AuthorID = nil
	})
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	wantTemplates := lo.Map(templates, func(t test_db.Template, _ int) domain.Template {
		return domain.Template{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt.Truncate(1 * time.Microsecond),
			UpdatedAt: lo.ToPtr(t.UpdatedAt.Truncate(1 * time.Microsecond)),
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
			in := domain.TemplateListDefaultIn{
				Page:    tt.page,
				Size:    tt.size,
				Filter:  domain.TemplateListDefaultFilter{},
				Sorting: nil,
			}

			got, err := repo.ListDefault(ctx, in)
			require.NoError(s.T(), err)
			require.Len(t, got, len(tt.want))
			require.Equal(t, tt.want, got)
		})
	}
}
