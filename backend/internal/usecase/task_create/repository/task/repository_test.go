package task_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_Insert() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.ProjectID = nil
		t.AuthorID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template version
	version := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userID
	})
	versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

	in := domain.TaskCreateIn{
		VersionID: versionID,
		CreatorID: userID,
		Payload: map[string]string{
			"test1": "123",
			"test2": "456.789",
			"test3": "text",
		},
	}

	err = repo.Insert(ctx, in)
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntityByColumn(s.C(), "task", "version_id", versionID))
	}()

	gotTasks, err := test_db.SelectEntitiesByColumn[test_db.Task](s.C(), "task", "version_id", []int64{versionID})
	require.NoError(s.T(), err)
	require.Len(s.T(), gotTasks, 1)

	got := gotTasks[0]

	want := test_db.Task{
		ID:        got.ID,
		VersionID: versionID,
		Status:    string(task_domain.StatusCreated),
		Payload:   []byte(`{"test1": "123", "test2": "456.789", "test3": "text"}`),
		ResultID:  nil,
		Error:     nil,
		CreatorID: userID,
		CreatedAt: got.CreatedAt,
		UpdatedAt: nil,
	}

	require.Equal(s.T(), want, got)
}
