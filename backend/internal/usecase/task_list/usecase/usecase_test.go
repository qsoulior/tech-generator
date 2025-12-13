package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	in := domain.TaskListIn{
		Page: 1,
		Size: 10,
		Filter: domain.TaskListFilter{
			UserID:     1,
			TemplateID: 2,
		},
	}

	var want domain.TaskListOut
	require.NoError(t, gofakeit.Struct(&want))
	want.TotalPages = (want.TotalTasks + in.Size - 1) / in.Size

	tests := []struct {
		name  string
		setup func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository)
	}{
		{
			name: "IsProjectAuthor",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{ProjectAuthorID: 1, TemplateAuthorID: 2}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
				taskRepo.EXPECT().List(ctx, in).Return(want.Tasks, nil)
				taskRepo.EXPECT().GetTotal(ctx, in).Return(want.TotalTasks, nil)
			},
		},
		{
			name: "IsTemplateAuthor",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{ProjectAuthorID: 2, TemplateAuthorID: 1}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
				taskRepo.EXPECT().List(ctx, in).Return(want.Tasks, nil)
				taskRepo.EXPECT().GetTotal(ctx, in).Return(want.TotalTasks, nil)
			},
		},
		{
			name: "IsReader",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{TemplateUsers: []domain.TemplateUser{{ID: 1, Role: user_domain.RoleRead}}}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
				taskRepo.EXPECT().List(ctx, in).Return(want.Tasks, nil)
				taskRepo.EXPECT().GetTotal(ctx, in).Return(want.TotalTasks, nil)
			},
		},
		{
			name: "IsWriter",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{TemplateUsers: []domain.TemplateUser{{ID: 1, Role: user_domain.RoleWrite}}}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
				taskRepo.EXPECT().List(ctx, in).Return(want.Tasks, nil)
				taskRepo.EXPECT().GetTotal(ctx, in).Return(want.TotalTasks, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			versionRepo := NewMockversionRepository(ctrl)
			taskRepo := NewMocktaskRepository(ctrl)
			tt.setup(versionRepo, taskRepo)

			usecase := New(versionRepo, taskRepo)

			got, err := usecase.Handle(ctx, in)
			require.NoError(t, err)
			require.Equal(t, want, *got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	in := domain.TaskListIn{
		Page: 1,
		Size: 10,
		Filter: domain.TaskListFilter{
			UserID:     1,
			TemplateID: 2,
		},
	}

	tests := []struct {
		name  string
		in    domain.TaskListIn
		setup func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository)
		want  string
	}{
		{
			name:  "in_Validate",
			in:    domain.TaskListIn{Page: 0, Size: 0},
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {},
			want:  domain.ErrValueInvalid.Error(),
		},
		{
			name: "versionRepo_GetByID",
			in:   in,
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrVersionNotFound",
			in:   in,
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(nil, nil)
			},
			want: domain.ErrVersionNotFound.Error(),
		},
		{
			name: "domain_ErrVersionInvalid",
			in:   in,
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{ProjectAuthorID: 2, TemplateAuthorID: 3}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
			},
			want: domain.ErrVersionInvalid.Error(),
		},
		{
			name: "taskRepo_List",
			in:   in,
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{ProjectAuthorID: 1, TemplateAuthorID: 2}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
				taskRepo.EXPECT().List(ctx, in).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "taskRepo_GetTotal",
			in:   in,
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository) {
				version := domain.Version{ProjectAuthorID: 1, TemplateAuthorID: 2}
				versionRepo.EXPECT().GetByID(ctx, in.Filter.TemplateID).Return(&version, nil)
				taskRepo.EXPECT().List(ctx, in).Return([]domain.Task{}, nil)
				taskRepo.EXPECT().GetTotal(ctx, in).Return(int64(0), errors.New("test3"))
			},
			want: "test3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			versionRepo := NewMockversionRepository(ctrl)
			taskRepo := NewMocktaskRepository(ctrl)
			tt.setup(versionRepo, taskRepo)

			usecase := New(versionRepo, taskRepo)

			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
