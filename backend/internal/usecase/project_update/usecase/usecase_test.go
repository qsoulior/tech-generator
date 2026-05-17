package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		in    domain.ProjectUpdateIn
		setup func(projectRepo *MockprojectRepository)
	}{
		{
			name: "IsAuthor",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
				projectRepo.EXPECT().UpdateByID(ctx, int64(10), "new").Return(nil)
			},
		},
		{
			name: "IsMaintainer",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 2, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{
					AuthorID: 1,
					Users: []domain.ProjectUser{
						{ID: 2, Role: user_domain.RoleMaintain},
					},
				}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
				projectRepo.EXPECT().UpdateByID(ctx, int64(10), "new").Return(nil)
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
		name    string
		in      domain.ProjectUpdateIn
		setup   func(projectRepo *MockprojectRepository)
		want    string
		wantVal bool
	}{
		{
			name:    "ValidationEmptyName",
			in:      domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: ""},
			setup:   func(_ *MockprojectRepository) {},
			want:    "name",
			wantVal: true,
		},
		{
			name: "projectRepo_GetByID",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid_NonMember",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 2}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "domain_ErrProjectInvalid_NonMaintainRole",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 3, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{
					AuthorID: 2,
					Users: []domain.ProjectUser{
						{ID: 3, Role: user_domain.RoleWrite},
					},
				}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "projectRepo_UpdateByID",
			in:   domain.ProjectUpdateIn{ProjectID: 10, UserID: 1, Name: "new"},
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 1}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
				projectRepo.EXPECT().UpdateByID(ctx, int64(10), "new").Return(errors.New("test2"))
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

			if tt.wantVal {
				var validationErr *error_domain.ValidationError
				require.ErrorAs(t, err, &validationErr)
			}
		})
	}
}
