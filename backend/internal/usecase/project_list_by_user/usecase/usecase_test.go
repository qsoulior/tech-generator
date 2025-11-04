package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	in := domain.ProjectListByUserIn{
		Page:   2,
		Size:   5,
		Filter: domain.ProjectListByUserFilter{UserID: 10},
		Sorting: &sorting_domain.Sorting{
			Attribute: "attr",
			Direction: "asc",
		},
	}

	want := domain.ProjectListByUserOut{
		Projects:      make([]domain.Project, 5),
		TotalProjects: 11,
		TotalPages:    3,
	}
	gofakeit.Slice(&want.Projects)

	projectRepo := NewMockprojectRepository(ctrl)
	projectRepo.EXPECT().ListByUserID(ctx, in).Return(want.Projects, nil)
	projectRepo.EXPECT().GetTotalByUserID(ctx, in).Return(want.TotalProjects, nil)

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
			name: "projectRepo_ListByAuthorID",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().ListByUserID(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			in:   domain.ProjectListByUserIn{Page: 1, Size: 1, Filter: domain.ProjectListByUserFilter{UserID: 10}},
			want: "test1",
		},
		{
			name: "projectRepo_GetTotalByUserID",
			setup: func(projectRepo *MockprojectRepository) {
				projectRepo.EXPECT().ListByUserID(ctx, gomock.Any()).Return([]domain.Project{}, nil)
				projectRepo.EXPECT().GetTotalByUserID(ctx, gomock.Any()).Return(int64(0), errors.New("test2"))
			},
			in:   domain.ProjectListByUserIn{Page: 1, Size: 1, Filter: domain.ProjectListByUserFilter{UserID: 10}},
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
