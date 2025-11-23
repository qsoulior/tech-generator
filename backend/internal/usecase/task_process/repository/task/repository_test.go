package task_repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
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

		// template
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.ProjectID = nil
			t.AuthorID = nil
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		// template version
		version := test_db.GenerateEntity(func(v *test_db.Version) {
			v.TemplateID = templateID
			v.AuthorID = &userID
		})
		versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

		want := domain.Task{
			VersionID: versionID,
			Payload: map[string]any{
				"test1": float64(123),
				"test2": 456.789,
				"test3": "text",
			},
		}

		payload, err := json.Marshal(want.Payload)
		require.NoError(t, err)

		// task
		task := test_db.GenerateEntity(func(t *test_db.Task) {
			t.CreatorID = userID
			t.VersionID = versionID
			t.ResultID = nil
			t.Payload = payload
			t.Error = nil
		})
		taskID, err := test_db.InsertEntityWithID[int64](s.C(), "task", task)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "task", taskID)) }()

		got, err := repo.GetByID(ctx, taskID)
		require.NoError(t, err)

		require.Equal(t, want, *got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func (s *repositorySuite) TestRepository_UpdateByID() {
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

	// task
	task := test_db.GenerateEntity(func(t *test_db.Task) {
		t.Status = string(task_domain.StatusCreated)
		t.CreatorID = userID
		t.VersionID = versionID
		t.ResultID = nil
		t.Payload = []byte("{\"test\":123}")
		t.Error = nil
	})
	taskID, err := test_db.InsertEntityWithID[int64](s.C(), "task", task)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "task", taskID)) }()

	// result
	result := test_db.GenerateEntity[test_db.Result]()
	resultID, err := test_db.InsertEntityWithID[int64](s.C(), "result", result)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "result", resultID)) }()

	// update #1
	taskUpdate := domain.TaskUpdate{
		ID:       taskID,
		Status:   task_domain.StatusSucceed,
		ResultID: &resultID,
		Error:    &task_domain.ProcessError{Message: "123"},
	}
	err = repo.UpdateByID(ctx, taskUpdate)
	require.NoError(s.T(), err)

	gotTasks, err := test_db.SelectEntitiesByID[test_db.Task](s.C(), "task", []int64{taskID})
	require.NoError(s.T(), err)
	require.Len(s.T(), gotTasks, 1)

	got := gotTasks[0]

	want := test_db.Task{
		ID:        taskID,
		VersionID: versionID,
		Status:    string(task_domain.StatusSucceed),
		Payload:   []byte("{\"test\": 123}"),
		ResultID:  &resultID,
		Error:     []byte("{\"message\": \"123\"}"),
		CreatorID: userID,
		CreatedAt: got.CreatedAt,
		UpdatedAt: got.UpdatedAt,
	}

	require.Equal(s.T(), want, got)

	// update #2
	taskUpdate = domain.TaskUpdate{
		ID:       taskID,
		Status:   task_domain.StatusFailed,
		ResultID: nil,
		Error:    nil,
	}
	err = repo.UpdateByID(ctx, taskUpdate)
	require.NoError(s.T(), err)

	gotTasks, err = test_db.SelectEntitiesByID[test_db.Task](s.C(), "task", []int64{taskID})
	require.NoError(s.T(), err)
	require.Len(s.T(), gotTasks, 1)

	got = gotTasks[0]

	want = test_db.Task{
		ID:        taskID,
		VersionID: versionID,
		Status:    string(task_domain.StatusFailed),
		Payload:   []byte("{\"test\": 123}"),
		ResultID:  nil,
		Error:     nil,
		CreatorID: userID,
		CreatedAt: got.CreatedAt,
		UpdatedAt: got.UpdatedAt,
	}

	require.Equal(s.T(), want, got)
}
