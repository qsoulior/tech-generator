package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	in := domain.ProjectUserListIn{UserID: 1, ProjectID: 10}
	want := []domain.ProjectUser{
		{ID: 2, Name: "alice", Email: "alice@example.com", Role: user_domain.RoleRead},
		{ID: 3, Name: "bob", Email: "bob@example.com", Role: user_domain.RoleWrite},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectRepo := NewMockprojectRepository(ctrl)
	projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&domain.Project{AuthorID: 1}, nil)

	projectUserRepo := NewMockprojectUserRepository(ctrl)
	projectUserRepo.EXPECT().GetByProjectID(ctx, int64(10)).Return(want, nil)

	usecase := New(projectRepo, projectUserRepo)
	got, err := usecase.Handle(ctx, in)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.ProjectUserListIn
		setup func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository)
		want  string
	}{
		{
			name: "projectRepo_GetByID",
			in:   domain.ProjectUserListIn{UserID: 1, ProjectID: 10},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			in:   domain.ProjectUserListIn{UserID: 1, ProjectID: 10},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid",
			in:   domain.ProjectUserListIn{UserID: 1, ProjectID: 10},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 9}, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "projectUserRepo_GetByProjectID",
			in:   domain.ProjectUserListIn{UserID: 1, ProjectID: 10},
			setup: func(projectRepo *MockprojectRepository, projectUserRepo *MockprojectUserRepository) {
				projectRepo.EXPECT().GetByID(ctx, gomock.Any()).Return(&domain.Project{AuthorID: 1}, nil)
				projectUserRepo.EXPECT().GetByProjectID(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			projectUserRepo := NewMockprojectUserRepository(ctrl)
			tt.setup(projectRepo, projectUserRepo)

			usecase := New(projectRepo, projectUserRepo)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
