package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.ProjectDeleteIn
		setup func(projectRepo *MockprojectRepository)
	}{
		{
			name: "IsAuthor",
			in:   domain.ProjectDeleteIn{ProjectID: 10, UserID: 1},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
				projectRepo.EXPECT().DeleteByID(ctx, int64(10)).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			tt.setup(projectRepo)

			usecase := New(projectRepo)
			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.ProjectDeleteIn
		setup func(projectRepo *MockprojectRepository)
		want  string
	}{
		{
			name: "projectRepo_GetByID",
			in:   domain.ProjectDeleteIn{ProjectID: 10, UserID: 1},
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			in:   domain.ProjectDeleteIn{ProjectID: 10, UserID: 1},
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			in:   domain.ProjectDeleteIn{ProjectID: 10, UserID: 1},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 2}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			in:   domain.ProjectDeleteIn{ProjectID: 10, UserID: 1},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
				projectRepo.EXPECT().DeleteByID(ctx, int64(10)).Return(errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			tt.setup(projectRepo)

			usecase := New(projectRepo)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
