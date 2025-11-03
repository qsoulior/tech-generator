package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectRepo := NewMockprojectRepository(ctrl)

	project := domain.Project{Name: "test", AuthorID: 1}
	projectRepo.EXPECT().Create(ctx, project).Return(nil)

	in := domain.ProjectCreateIn{Name: "test", AuthorID: 1}
	usecase := New(projectRepo)
	err := usecase.Handle(ctx, in)
	require.NoError(t, err)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		setup func(projectRepo *MockprojectRepository)
		in    domain.ProjectCreateIn
		want  string
	}{
		{
			name:  "in_Validate",
			setup: func(projectRepo *MockprojectRepository) {},
			in:    domain.ProjectCreateIn{Name: ""},
			want:  domain.ErrValueEmpty.Error(),
		},
		{
			name: "projectRepo_Create",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("test1"))
			},
			in:   domain.ProjectCreateIn{Name: "test", AuthorID: 1},
			want: "test1",
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
