package task_repository

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_List() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	users := test_db.GenerateEntities[test_db.User](2)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.AuthorID = &userIDs[0]
		t.ProjectID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// version
	version := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userIDs[0]
	})
	versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

	// tasks
	tasks := test_db.GenerateEntities(2, func(t *test_db.Task, i int) {
		t.VersionID = versionID
		t.CreatorID = userIDs[i]
		t.ResultID = nil
		t.Payload = []byte("{}")
		t.Error = nil
	})
	taskIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "task", tasks)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "task", taskIDs)) }()

	wantTasks := lo.Map(tasks, func(t test_db.Task, i int) domain.Task {
		return domain.Task{
			ID:            t.ID,
			VersionNumber: version.Number,
			Status:        task_domain.Status(t.Status),
			CreatorName:   users[i].Name,
			CreatedAt:     t.CreatedAt.Truncate(1 * time.Microsecond),
			UpdatedAt:     lo.ToPtr(t.UpdatedAt.Truncate(1 * time.Microsecond)),
		}
	})
	slices.SortFunc(wantTasks, func(a, b domain.Task) int { return int(b.ID - a.ID) })

	in := domain.TaskListIn{
		Page: 1,
		Size: int64(len(tasks)),
		Filter: domain.TaskListFilter{
			TemplateID: templateID,
		},
		Sorting: nil,
	}

	gotTasks, err := repo.List(ctx, in)
	require.NoError(s.T(), err)
	require.Equal(s.T(), wantTasks, gotTasks)

	gotTotal, err := repo.GetTotal(ctx, in)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), len(gotTasks), gotTotal)
}

func (s *repositorySuite) TestRepository_List_Filter() {
	ctx := context.Background()
	repo := New(s.C().DB())

	const N = 5

	// users
	users := test_db.GenerateEntities[test_db.User](N)
	userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

	randomUser := lo.Sample(users)

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.AuthorID = &userIDs[0]
		t.ProjectID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// version
	version := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userIDs[0]
	})
	versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

	// tasks
	tasks := test_db.GenerateEntities(N, func(t *test_db.Task, i int) {
		t.VersionID = versionID
		t.CreatorID = userIDs[i]
		t.ResultID = nil
		t.Payload = []byte("{}")
		t.Error = nil
	})
	taskIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "task", tasks)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "task", taskIDs)) }()

	tests := []struct {
		name       string
		filter     domain.TaskListFilter
		filterPred func(t domain.Task) bool
	}{
		{
			name: "CreatorID",
			filter: domain.TaskListFilter{
				TemplateID: templateID,
				CreatorID:  &randomUser.ID,
			},
			filterPred: func(t domain.Task) bool {
				return t.CreatorName == randomUser.Name
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TaskListIn{
				Page:    1,
				Size:    int64(len(tasks)),
				Filter:  tt.filter,
				Sorting: nil,
			}

			gotTasks, err := repo.List(ctx, in)
			require.NoError(s.T(), err)

			isFiltered := lo.EveryBy(gotTasks, tt.filterPred)
			require.True(t, isFiltered)

			gotTotal, err := repo.GetTotal(ctx, in)
			require.NoError(s.T(), err)
			require.EqualValues(t, len(gotTasks), gotTotal)
		})
	}
}

func (s *repositorySuite) TestRepository_List_Sorting() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.AuthorID = &userID
		t.ProjectID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// version
	version := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userID
	})
	versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

	// tasks
	tasks := test_db.GenerateEntities(2, func(t *test_db.Task, i int) {
		t.VersionID = versionID
		t.CreatorID = userID
		t.ResultID = nil
		t.Payload = []byte("{}")
		t.Error = nil
	})
	taskIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "task", tasks)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "task", taskIDs)) }()

	tests := []struct {
		name       string
		sorting    *sorting_domain.Sorting
		comparator func(a, b domain.Task) int
	}{
		{
			name:    "Default",
			sorting: nil,
			comparator: func(a, b domain.Task) int {
				return int(b.ID - a.ID)
			},
		},
		{
			name: "Invalid",
			sorting: &sorting_domain.Sorting{
				Attribute: "invalid",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Task) int {
				return int(b.ID - a.ID)
			},
		},
		{
			name: "CreatedAt_ASC",
			sorting: &sorting_domain.Sorting{
				Attribute: "created_at",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Task) int {
				return a.CreatedAt.Compare(b.CreatedAt)
			},
		},
		{
			name: "CreatedAt_DESC",
			sorting: &sorting_domain.Sorting{
				Attribute: "created_at",
				Direction: "DESC",
			},
			comparator: func(a, b domain.Task) int {
				return b.CreatedAt.Compare(a.CreatedAt)
			},
		},
		{
			name: "UpdatedAt_ASC",
			sorting: &sorting_domain.Sorting{
				Attribute: "updated_at",
				Direction: "ASC",
			},
			comparator: func(a, b domain.Task) int {
				return a.UpdatedAt.Compare(*b.UpdatedAt)
			},
		},
		{
			name: "UpdatedAt_DESC",
			sorting: &sorting_domain.Sorting{
				Attribute: "updated_at",
				Direction: "DESC",
			},
			comparator: func(a, b domain.Task) int {
				return b.UpdatedAt.Compare(*a.UpdatedAt)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TaskListIn{
				Page: 1,
				Size: int64(len(tasks)),
				Filter: domain.TaskListFilter{
					TemplateID: templateID,
				},
				Sorting: tt.sorting,
			}

			got, err := repo.List(ctx, in)
			require.NoError(s.T(), err)

			isSorted := slices.IsSortedFunc(got, tt.comparator)
			require.True(t, isSorted)
		})
	}
}

func (s *repositorySuite) TestRepository_List_Pagination() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// users
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.AuthorID = &userID
		t.ProjectID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// version
	version := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userID
	})
	versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

	// tasks
	tasks := test_db.GenerateEntities(5, func(t *test_db.Task, i int) {
		t.VersionID = versionID
		t.CreatorID = userID
		t.ResultID = nil
		t.Payload = []byte("{}")
		t.Error = nil
	})
	taskIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "task", tasks)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "task", taskIDs)) }()

	slices.SortFunc(taskIDs, func(a, b int64) int { return int(b - a) })

	tests := []struct {
		name string
		page int64
		size int64
		want []int64
	}{
		{
			name: "None",
			page: 2,
			size: int64(len(tasks)),
			want: []int64{},
		},
		{
			name: "One",
			page: 1,
			size: 1,
			want: taskIDs[:1],
		},
		{
			name: "Many",
			page: 2,
			size: 2,
			want: taskIDs[2:4],
		},
		{
			name: "All",
			page: 1,
			size: int64(len(tasks)),
			want: taskIDs,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			in := domain.TaskListIn{
				Page:    tt.page,
				Size:    tt.size,
				Filter:  domain.TaskListFilter{TemplateID: templateID},
				Sorting: nil,
			}

			got, err := repo.List(ctx, in)
			require.NoError(s.T(), err)
			require.Len(t, got, len(tt.want))

			gotIDs := lo.Map(got, func(t domain.Task, _ int) int64 { return t.ID })
			require.Equal(t, tt.want, gotIDs)
		})
	}
}
