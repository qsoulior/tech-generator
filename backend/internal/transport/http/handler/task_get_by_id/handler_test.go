package task_get_by_id_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
)

func TestHandler_TaskGetByID_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TaskGetByIDParams{TaskID: 9, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createdAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	taskErr := task_domain.ProcessError{
		Message: "fail",
		VariableErrors: []task_domain.VariableError{
			{
				ID:      11,
				Name:    "v1",
				Value:   "42",
				Message: "bad",
				ConstraintErrors: []task_domain.ConstraintError{
					{ID: 21, Name: "c1", Expression: "v1 > 100", Message: "broken"},
				},
			},
		},
	}
	out := &domain.TaskGetByIDOut{
		Task: domain.Task{
			ID:          9,
			VersionID:   7,
			Status:      task_domain.StatusFailed,
			Payload:     map[string]string{"k": "v"},
			Error:       &taskErr,
			CreatorName: "alice",
			CreatedAt:   createdAt,
			UpdatedAt:   &updatedAt,
		},
		Result: []byte("payload"),
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, domain.TaskGetByIDIn{TaskID: 9, UserID: 1}).Return(out, nil)

	handler := New(usecase)
	got, err := handler.TaskGetByID(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TaskGetByIDResponse)
	require.True(t, ok, "expected *api.TaskGetByIDResponse, got %T", got)
	require.Equal(t, int64(9), resp.Task.ID)
	require.Equal(t, int64(7), resp.Task.VersionID)
	require.Equal(t, api.TaskStatus(task_domain.StatusFailed), resp.Task.Status)
	require.Equal(t, "alice", resp.Task.CreatorName)
	require.Equal(t, createdAt, resp.Task.CreatedAt)

	gotUpdatedAt, ok := resp.Task.UpdatedAt.Get()
	require.True(t, ok)
	require.Equal(t, updatedAt, gotUpdatedAt)

	gotErr, ok := resp.Task.Error.Get()
	require.True(t, ok)
	gotMsg, ok := gotErr.Message.Get()
	require.True(t, ok)
	require.Equal(t, "fail", gotMsg)
	require.Len(t, gotErr.VariableErrors, 1)
	require.Equal(t, int64(11), gotErr.VariableErrors[0].ID)
	require.Equal(t, "v1", gotErr.VariableErrors[0].Name)
	gotValue, ok := gotErr.VariableErrors[0].Value.Get()
	require.True(t, ok)
	require.Equal(t, "42", gotValue)
	require.Len(t, gotErr.VariableErrors[0].ConstraintErrors, 1)
	require.Equal(t, int64(21), gotErr.VariableErrors[0].ConstraintErrors[0].ID)
	require.Equal(t, "v1 > 100", gotErr.VariableErrors[0].ConstraintErrors[0].Expression)
	require.Equal(t, []byte("payload"), resp.Result)
}

func TestHandler_TaskGetByID_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.TaskGetByIDParams{TaskID: 9, XUserID: 1}

	tests := []struct {
		name string
		err  error
	}{
		{name: "NotFound", err: domain.ErrTaskNotFound},
		{name: "Invalid", err: domain.ErrTaskInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.TaskGetByID(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TaskGetByID_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TaskGetByIDParams{TaskID: 9, XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TaskGetByID(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "task get by id usecase")
	require.ErrorContains(t, err, "boom")
}
