package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var in domain.ProjectListByUserIn
	require.NoError(t, gofakeit.Struct(&in))

	var want domain.ProjectListByUserOut
	require.NoError(t, gofakeit.Struct(&want))

	projectRepo := NewMockprojectRepository(ctrl)
	projectRepo.EXPECT().ListByAuthorID(ctx, in).Return(want.Owned, nil)
	projectRepo.EXPECT().ListByProjectUserID(ctx, in).Return(want.Shared, nil)

	usecase := New(projectRepo)
	got, err := usecase.Handle(ctx, in)
	require.NoError(t, err)
	require.Equal(t, want, *got)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		setup func(projectRepo *MockprojectRepository)
		in    domain.ProjectListByUserIn
		want  string
	}{
		{
			name:  "in_Validate",
			setup: func(projectRepo *MockprojectRepository) {},
			in:    domain.ProjectListByUserIn{Page: 0, Size: 0},
			want:  domain.ErrValueInvalid.Error(),
		},
		{
			name: "",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().ListByAuthorID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			in:   domain.ProjectListByUserIn{UserID: 10, Page: 1, Size: 1},
			want: "test1",
		},
		{
			name: "",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().ListByAuthorID(ctx, gomock.Any()).Return([]domain.Project{}, nil)
				projectRepo.EXPECT().ListByProjectUserID(ctx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			in:   domain.ProjectListByUserIn{UserID: 10, Page: 1, Size: 1},
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
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
