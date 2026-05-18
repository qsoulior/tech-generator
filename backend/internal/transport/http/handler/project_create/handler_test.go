package project_create_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

func TestHandler_ProjectCreate_Success(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectCreateRequest{Name: "p"}
	params := api.ProjectCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectCreateIn{Name: "p", AuthorID: 1}).
		Return(nil)

	handler := New(usecase)
	got, err := handler.ProjectCreate(ctx, req, params)
	require.NoError(t, err)

	_, ok := got.(*api.ProjectCreateCreated)
	require.True(t, ok, "expected *api.ProjectCreateCreated, got %T", got)
}

func TestHandler_ProjectCreate_ValidationError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectCreateRequest{Name: ""}
	params := api.ProjectCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validationErr := error_domain.NewValidationError("name", errors.New("empty"))

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectCreateIn{Name: "", AuthorID: 1}).
		Return(validationErr)

	handler := New(usecase)
	got, err := handler.ProjectCreate(ctx, req, params)
	require.NoError(t, err)

	resp, ok := got.(*api.Error)
	require.True(t, ok, "expected *api.Error, got %T", got)
	require.Equal(t, validationErr.Error(), resp.Message)
}

func TestHandler_ProjectCreate_InternalError(t *testing.T) {
	ctx := context.Background()
	req := &api.ProjectCreateRequest{Name: "p"}
	params := api.ProjectCreateParams{XUserID: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().
		Handle(ctx, domain.ProjectCreateIn{Name: "p", AuthorID: 1}).
		Return(errors.New("boom"))

	handler := New(usecase)
	got, err := handler.ProjectCreate(ctx, req, params)
	require.Nil(t, got)
	require.ErrorContains(t, err, "project create usecase")
	require.ErrorContains(t, err, "boom")
}
