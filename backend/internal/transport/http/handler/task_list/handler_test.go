package task_list_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

func TestHandler_TaskList_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TaskListParams{XUserID: 1, Page: 1, Size: 10, TemplateID: 3}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createdAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	in := domain.TaskListIn{
		Page:   1,
		Size:   10,
		Filter: domain.TaskListFilter{UserID: 1, TemplateID: 3},
	}
	out := &domain.TaskListOut{
		Tasks: []domain.Task{{
			ID:            9,
			VersionNumber: 2,
			Status:        task_domain.StatusSucceed,
			CreatorName:   "alice",
			CreatedAt:     createdAt,
			UpdatedAt:     &updatedAt,
		}},
		TotalTasks: 1,
		TotalPages: 1,
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(out, nil)

	handler := New(usecase)
	got, err := handler.TaskList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TaskListResponse)
	require.True(t, ok, "expected *api.TaskListResponse, got %T", got)
	require.Equal(t, int64(1), resp.TotalTasks)
	require.Equal(t, int64(1), resp.TotalPages)
	require.Len(t, resp.Tasks, 1)
	require.Equal(t, int64(9), resp.Tasks[0].ID)
	require.Equal(t, int64(2), resp.Tasks[0].VersionNumber)
	require.Equal(t, api.TaskStatus(task_domain.StatusSucceed), resp.Tasks[0].Status)
	gotUpdatedAt, ok := resp.Tasks[0].UpdatedAt.Get()
	require.True(t, ok)
	require.Equal(t, updatedAt, gotUpdatedAt)
}

func TestHandler_TaskList_PassesFilterAndSorting(t *testing.T) {
	ctx := context.Background()
	params := api.TaskListParams{XUserID: 1, Page: 2, Size: 5, TemplateID: 3}
	params.CreatorID.SetTo(42)
	params.Sorting.SetTo(api.Sorting{Attribute: "createdAt", Direction: api.SortingDirectionDESC})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	creatorID := int64(42)
	in := domain.TaskListIn{
		Page: 2,
		Size: 5,
		Filter: domain.TaskListFilter{
			UserID:     1,
			TemplateID: 3,
			CreatorID:  &creatorID,
		},
		Sorting: &sorting_domain.Sorting{
			Attribute: "createdAt",
			Direction: sorting_domain.SortingDirectionDesc,
		},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&domain.TaskListOut{}, nil)

	handler := New(usecase)
	_, err := handler.TaskList(ctx, params)
	require.NoError(t, err)
}

func TestHandler_TaskList_BaseError(t *testing.T) {
	ctx := context.Background()
	params := api.TaskListParams{XUserID: 1, Page: 1, Size: 10, TemplateID: 3}

	tests := []struct {
		name string
		err  error
	}{
		{name: "TemplateNotFound", err: domain.ErrTemplateNotFound},
		{name: "TemplateInvalid", err: domain.ErrTemplateInvalid},
		{name: "WrappedBaseError", err: error_domain.NewBaseError("custom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, tt.err)

			handler := New(usecase)
			got, err := handler.TaskList(ctx, params)
			require.NoError(t, err)

			resp, ok := got.(*api.Error)
			require.True(t, ok, "expected *api.Error, got %T", got)
			require.Equal(t, tt.err.Error(), resp.Message)
		})
	}
}

func TestHandler_TaskList_ValidationError(t *testing.T) {
	ctx := context.Background()
	params := api.TaskListParams{XUserID: 1, Page: 0, Size: 10, TemplateID: 3}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("page", errors.New("invalid"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, validationErr)

	handler := New(usecase)
	got, err := handler.TaskList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_TaskList_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TaskListParams{XUserID: 1, Page: 1, Size: 10, TemplateID: 3}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TaskList(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "task list usecase")
	require.ErrorContains(t, err, "boom")
}
