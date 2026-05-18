package project_list_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

func TestHandler_ProjectList_Success(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectListParams{XUserID: 1, Page: 1, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	in := domain.ProjectListByUserIn{
		Page:   1,
		Size:   10,
		Filter: domain.ProjectListByUserFilter{UserID: 1},
	}
	out := domain.ProjectListByUserOut{
		Projects: []domain.Project{
			{ID: 7, Name: "p", AuthorName: "alice"},
		},
		TotalProjects: 1,
		TotalPages:    1,
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&out, nil)

	handler := New(usecase)
	got, err := handler.ProjectList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.ProjectListResponse)
	require.True(t, ok, "expected *api.ProjectListResponse, got %T", got)
	require.Equal(t, int64(1), resp.TotalProjects)
	require.Equal(t, int64(1), resp.TotalPages)
	require.Len(t, resp.Projects, 1)
	require.Equal(t, int64(7), resp.Projects[0].ID)
	require.Equal(t, "p", resp.Projects[0].Name)
	require.Equal(t, "alice", resp.Projects[0].AuthorName)
}

func TestHandler_ProjectList_PassesFilterAndSorting(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectListParams{XUserID: 1, Page: 2, Size: 5}
	params.ProjectName.SetTo("foo")
	params.Sorting.SetTo(api.Sorting{Attribute: "name", Direction: api.SortingDirectionDESC})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectName := "foo"
	in := domain.ProjectListByUserIn{
		Page: 2,
		Size: 5,
		Filter: domain.ProjectListByUserFilter{
			UserID:      1,
			ProjectName: &projectName,
		},
		Sorting: &sorting_domain.Sorting{
			Attribute: "name",
			Direction: sorting_domain.SortingDirectionDesc,
		},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&domain.ProjectListByUserOut{}, nil)

	handler := New(usecase)
	_, err := handler.ProjectList(ctx, params)
	require.NoError(t, err)
}

func TestHandler_ProjectList_ValidationError(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectListParams{XUserID: 1, Page: 0, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("page", errors.New("invalid"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, validationErr)

	handler := New(usecase)
	got, err := handler.ProjectList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_ProjectList_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.ProjectListParams{XUserID: 1, Page: 1, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.ProjectList(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "project list by user usecase")
	require.ErrorContains(t, err, "boom")
}
