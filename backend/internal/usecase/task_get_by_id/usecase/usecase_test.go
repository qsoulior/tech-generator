package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	var in domain.TaskGetByIDIn
	require.NoError(t, gofakeit.Struct(&in))

	tests := []struct {
		name  string
		setup func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository)
		want  *domain.TaskGetByIDOut
	}{
		{
			name: "StatusSucceed/IsProjectAuthor",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusSucceed,
					Payload:     map[string]any{"test": 1},
					ResultID:    lo.ToPtr[int64](3),
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				version := domain.Version{ProjectAuthorID: in.UserID}
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(&version, nil)

				result := []byte{1, 2, 3}
				resultRepo.EXPECT().GetDataByID(ctx, *task.ResultID).Return(result, nil)
			},
			want: &domain.TaskGetByIDOut{
				Task: domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusSucceed,
					Payload:     map[string]any{"test": 1},
					ResultID:    lo.ToPtr[int64](3),
					CreatorName: "test",
				},
				Result: []byte{1, 2, 3},
			},
		},
		{
			name: "StatusFailed/IsTemplateAuthor",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusFailed,
					Payload:     map[string]any{"test": 1},
					Error:       &task_domain.ProcessError{Message: "test"},
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				version := domain.Version{TemplateAuthorID: in.UserID}
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(&version, nil)
			},
			want: &domain.TaskGetByIDOut{
				Task: domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusFailed,
					Payload:     map[string]any{"test": 1},
					Error:       &task_domain.ProcessError{Message: "test"},
					CreatorName: "test",
				},
			},
		},
		{
			name: "StatusFailed/IsReader",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusFailed,
					Error:       &task_domain.ProcessError{Message: "test"},
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				version := domain.Version{TemplateUsers: []domain.TemplateUser{{ID: in.UserID, Role: user_domain.RoleRead}}}
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(&version, nil)
			},
			want: &domain.TaskGetByIDOut{
				Task: domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusFailed,
					Error:       &task_domain.ProcessError{Message: "test"},
					CreatorName: "test",
				},
			},
		},
		{
			name: "StatusFailed/IsWriter",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusFailed,
					Error:       &task_domain.ProcessError{Message: "test"},
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				version := domain.Version{TemplateUsers: []domain.TemplateUser{{ID: in.UserID, Role: user_domain.RoleWrite}}}
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(&version, nil)
			},
			want: &domain.TaskGetByIDOut{
				Task: domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusFailed,
					Error:       &task_domain.ProcessError{Message: "test"},
					CreatorName: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepo := NewMocktaskRepository(ctrl)
			versionRepo := NewMockversionRepository(ctrl)
			resultRepo := NewMockresultRepository(ctrl)

			tt.setup(taskRepo, versionRepo, resultRepo)

			usecase := New(taskRepo, versionRepo, resultRepo)

			got, err := usecase.Handle(ctx, in)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	var in domain.TaskGetByIDIn
	require.NoError(t, gofakeit.Struct(&in))

	tests := []struct {
		name  string
		setup func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository)
		want  string
	}{
		{
			name: "taskRepo_GetByID",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrNotFound_#1",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(nil, nil)
			},
			want: domain.ErrTaskNotFound.Error(),
		},
		{
			name: "domain_ErrNotFound_#2",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusSucceed,
					ResultID:    lo.ToPtr[int64](3),
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(nil, nil)
			},
			want: domain.ErrTaskNotFound.Error(),
		},
		{
			name: "domain_ErrTaskInvalid",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusSucceed,
					ResultID:    lo.ToPtr[int64](3),
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				version := domain.Version{}
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(&version, nil)
			},
			want: domain.ErrTaskInvalid.Error(),
		},
		{
			name: "versionRepo_GetByID",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusSucceed,
					ResultID:    lo.ToPtr[int64](3),
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "resultRepo_GetDataByID",
			setup: func(taskRepo *MocktaskRepository, versionRepo *MockversionRepository, resultRepo *MockresultRepository) {
				task := domain.Task{
					ID:          1,
					VersionID:   2,
					Status:      task_domain.StatusSucceed,
					ResultID:    lo.ToPtr[int64](3),
					CreatorName: "test",
				}
				taskRepo.EXPECT().GetByID(ctx, in.TaskID).Return(&task, nil)

				version := domain.Version{ProjectAuthorID: in.UserID}
				versionRepo.EXPECT().GetByID(ctx, task.VersionID).Return(&version, nil)

				resultRepo.EXPECT().GetDataByID(ctx, *task.ResultID).Return(nil, errors.New("test3"))
			},
			want: "test3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepo := NewMocktaskRepository(ctrl)
			versionRepo := NewMockversionRepository(ctrl)
			resultRepo := NewMockresultRepository(ctrl)

			tt.setup(taskRepo, versionRepo, resultRepo)

			usecase := New(taskRepo, versionRepo, resultRepo)

			_, err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
