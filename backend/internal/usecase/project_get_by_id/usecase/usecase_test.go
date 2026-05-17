package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	in := domain.ProjectGetByIDIn{ProjectID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(projectRepo *MockprojectRepository)
		want  domain.ProjectGetByIDOut
	}{
		{
			name: "IsAuthor",
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{
					Name:       "test",
					AuthorID:   1,
					AuthorName: "author",
				}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ProjectGetByIDOut{Name: "test", AuthorName: "author"},
		},
		{
			name: "IsReader",
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{
					Name:       "test",
					AuthorID:   2,
					AuthorName: "author",
					Users:      []domain.ProjectUser{{ID: 1, Role: user_domain.RoleRead}},
				}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ProjectGetByIDOut{Name: "test", AuthorName: "author"},
		},
		{
			name: "IsWriter",
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{
					Name:       "test",
					AuthorID:   2,
					AuthorName: "author",
					Users:      []domain.ProjectUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ProjectGetByIDOut{Name: "test", AuthorName: "author"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			tt.setup(projectRepo)

			usecase := New(projectRepo)

			got, err := usecase.Handle(ctx, in)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	in := domain.ProjectGetByIDIn{ProjectID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(projectRepo *MockprojectRepository)
		want  string
	}{
		{
			name: "projectRepo_GetByID",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrProjectNotFound",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrProjectNotFound.Error(),
		},
		{
			name: "domain_ErrProjectInvalid/Stranger",
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{AuthorID: 2}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
		{
			name: "domain_ErrProjectInvalid/OtherUser",
			setup: func(projectRepo *MockprojectRepository) {
				project := domain.Project{
					AuthorID: 2,
					Users:    []domain.ProjectUser{{ID: 3, Role: user_domain.RoleWrite}},
				}
				projectRepo.EXPECT().GetByID(ctx, int64(10)).Return(&project, nil)
			},
			want: domain.ErrProjectInvalid.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projectRepo := NewMockprojectRepository(ctrl)
			tt.setup(projectRepo)

			usecase := New(projectRepo)

			_, err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
