package task_repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
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

		// result
		result := test_db.GenerateEntity[test_db.Result]()
		resultID, err := test_db.InsertEntityWithID[int64](s.C(), "result", result)
		require.NoError(s.T(), err)
		defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "result", resultID)) }()

		want := domain.Task{
			ID:        gofakeit.Int64(),
			VersionID: versionID,
			Status:    task_domain.StatusInProgress,
			Payload: map[string]any{
				"test1": float64(123),
				"test2": 456.789,
				"test3": "text",
			},
			ResultID: &resultID,
			Error: &task_domain.ProcessError{
				Message: "msg1",
				VariableErrors: []task_domain.VariableError{
					{
						ID:      2,
						Name:    "name2",
						Message: "msg2",
						ConstraintErrors: []task_domain.ConstraintError{
							{
								ID:      3,
								Name:    "name3",
								Message: "msg3",
							},
						},
					},
				},
			},
			CreatorName: user.Name,
			CreatedAt:   gofakeit.Date().Truncate(1 * time.Microsecond),
			UpdatedAt:   lo.ToPtr(gofakeit.Date().Truncate(1 * time.Microsecond)),
		}

		payload, err := json.Marshal(want.Payload)
		require.NoError(t, err)

		taskError, err := json.Marshal(want.Error)
		require.NoError(t, err)

		// task
		task := test_db.Task{
			ID:        want.ID,
			VersionID: versionID,
			Status:    string(want.Status),
			Payload:   payload,
			ResultID:  &resultID,
			Error:     taskError,
			CreatorID: userID,
			CreatedAt: want.CreatedAt,
			UpdatedAt: want.UpdatedAt,
		}
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
