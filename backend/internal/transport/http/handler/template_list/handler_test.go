package template_list_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

func TestHandler_TemplateList_Success(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateListParams{XUserID: 1, ProjectID: 3, Page: 1, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createdAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	in := domain.TemplateListByUserIn{
		Page:   1,
		Size:   10,
		Filter: domain.TemplateListByUserFilter{UserID: 1, ProjectID: 3},
	}
	out := &domain.TemplateListByUserOut{
		Templates: []domain.Template{{
			ID:         5,
			Name:       "tmpl",
			AuthorName: "alice",
			CreatedAt:  createdAt,
			UpdatedAt:  &updatedAt,
		}},
		TotalTemplates: 1,
		TotalPages:     1,
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(out, nil)

	handler := New(usecase)
	got, err := handler.TemplateList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.TemplateListResponse)
	require.True(t, ok, "expected *api.TemplateListResponse, got %T", got)
	require.Equal(t, int64(1), resp.TotalTemplates)
	require.Equal(t, int64(1), resp.TotalPages)
	require.Len(t, resp.Templates, 1)
	require.Equal(t, int64(5), resp.Templates[0].ID)
	require.Equal(t, "tmpl", resp.Templates[0].Name)
	require.Equal(t, "alice", resp.Templates[0].AuthorName)
	require.Equal(t, createdAt, resp.Templates[0].CreatedAt)
	gotUpdatedAt, ok := resp.Templates[0].UpdatedAt.Get()
	require.True(t, ok)
	require.Equal(t, updatedAt, gotUpdatedAt)
}

func TestHandler_TemplateList_PassesFilterAndSorting(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateListParams{XUserID: 1, ProjectID: 3, Page: 2, Size: 5}
	params.TemplateName.SetTo("foo")
	params.Sorting.SetTo(api.Sorting{Attribute: "name", Direction: api.SortingDirectionASC})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "foo"
	in := domain.TemplateListByUserIn{
		Page: 2,
		Size: 5,
		Filter: domain.TemplateListByUserFilter{
			UserID:       1,
			ProjectID:    3,
			TemplateName: &templateName,
		},
		Sorting: &sorting_domain.Sorting{
			Attribute: "name",
			Direction: sorting_domain.SortingDirectionAsc,
		},
	}

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, in).Return(&domain.TemplateListByUserOut{}, nil)

	handler := New(usecase)
	_, err := handler.TemplateList(ctx, params)
	require.NoError(t, err)
}

func TestHandler_TemplateList_ValidationError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateListParams{XUserID: 1, ProjectID: 3, Page: 0, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("page", errors.New("invalid"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, validationErr)

	handler := New(usecase)
	got, err := handler.TemplateList(ctx, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_TemplateList_InternalError(t *testing.T) {
	ctx := context.Background()
	params := api.TemplateListParams{XUserID: 1, ProjectID: 3, Page: 1, Size: 10}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("boom"))

	handler := New(usecase)
	got, err := handler.TemplateList(ctx, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "template list by user usecase")
	require.ErrorContains(t, err, "boom")
}
