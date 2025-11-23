package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	taskID := gofakeit.Int64()

	tests := []struct {
		name  string
		setup func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository)
	}{
		{
			name: "Success",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				var task domain.Task
				_ = gofakeit.Struct(&task)

				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&task, nil)

				taskUpdate := domain.TaskUpdate{ID: taskID, Status: task_domain.StatusInProgress}
				taskRepo.EXPECT().UpdateByID(ctx, taskUpdate).Return(nil)

				var version domain.Version
				_ = gofakeit.Struct(&version)
				versionGetService.EXPECT().Handle(ctx, task.VersionID).Return(&version, nil)

				variableProcessIn := domain.VariableProcessIn{Variables: version.Variables, Payload: task.Payload}
				variableValues := gofakeit.Map()
				variableProcessService.EXPECT().Handle(ctx, variableProcessIn).Return(variableValues, nil)

				dataProcessIn := domain.DataProcessIn{Values: variableValues, Data: version.Data}
				result := []byte{1, 2, 3}
				dataProcessService.EXPECT().Handle(ctx, dataProcessIn).Return(result, nil)

				resultID := gofakeit.Int64()
				resultRepo.EXPECT().Insert(ctx, result).Return(resultID, nil)

				taskUpdate = domain.TaskUpdate{ID: taskID, Status: task_domain.StatusSucceed, ResultID: &resultID}
				taskRepo.EXPECT().UpdateByID(ctx, taskUpdate).Return(nil)
			},
		},
		{
			name: "variableProcessService_ProcessError",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				var task domain.Task
				_ = gofakeit.Struct(&task)

				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&task, nil)

				taskUpdate := domain.TaskUpdate{ID: taskID, Status: task_domain.StatusInProgress}
				taskRepo.EXPECT().UpdateByID(ctx, taskUpdate).Return(nil)

				var version domain.Version
				_ = gofakeit.Struct(&version)
				versionGetService.EXPECT().Handle(ctx, task.VersionID).Return(&version, nil)

				variableProcessIn := domain.VariableProcessIn{Variables: version.Variables, Payload: task.Payload}
				err := &domain.ProcessError{Message: "test1"}
				variableProcessService.EXPECT().Handle(ctx, variableProcessIn).Return(nil, err)

				taskUpdate = domain.TaskUpdate{ID: taskID, Status: task_domain.StatusFailed, Error: err}
				taskRepo.EXPECT().UpdateByID(ctx, taskUpdate).Return(nil)
			},
		},
		{
			name: "dataProcessService_ProcessError",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				var task domain.Task
				_ = gofakeit.Struct(&task)

				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&task, nil)

				taskUpdate := domain.TaskUpdate{ID: taskID, Status: task_domain.StatusInProgress}
				taskRepo.EXPECT().UpdateByID(ctx, taskUpdate).Return(nil)

				var version domain.Version
				_ = gofakeit.Struct(&version)
				versionGetService.EXPECT().Handle(ctx, task.VersionID).Return(&version, nil)

				variableProcessIn := domain.VariableProcessIn{Variables: version.Variables, Payload: task.Payload}
				variableValues := gofakeit.Map()
				variableProcessService.EXPECT().Handle(ctx, variableProcessIn).Return(variableValues, nil)

				err := &domain.ProcessError{Message: "test2"}
				dataProcessIn := domain.DataProcessIn{Values: variableValues, Data: version.Data}
				dataProcessService.EXPECT().Handle(ctx, dataProcessIn).Return(nil, err)

				taskUpdate = domain.TaskUpdate{ID: taskID, Status: task_domain.StatusFailed, Error: err}
				taskRepo.EXPECT().UpdateByID(ctx, taskUpdate).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepo := NewMocktaskRepository(ctrl)
			versionGetService := NewMockversionGetService(ctrl)
			variableProcessService := NewMockvariableProcessService(ctrl)
			dataProcessService := NewMockdataProcessService(ctrl)
			resultRepo := NewMockresultRepository(ctrl)

			tt.setup(taskRepo, versionGetService, variableProcessService, dataProcessService, resultRepo)

			usecase := New(taskRepo, versionGetService, variableProcessService, dataProcessService, resultRepo)
			err := usecase.Handle(ctx, domain.TaskProcessIn{TaskID: taskID})
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	taskID := gofakeit.Int64()

	tests := []struct {
		name  string
		setup func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository)
		want  string
	}{
		{
			name: "taskRepo_GetByID",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTaskNotFound",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(nil, nil)
			},
			want: domain.ErrTaskNotFound.Error(),
		},
		{
			name: "taskRepo_UpdateByID_#1",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "versionGetService_Error",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(nil)
				versionGetService.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("test3"))
			},
			want: "test3",
		},
		{
			name: "variableProcessService_Error",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(nil)
				versionGetService.EXPECT().Handle(ctx, gomock.Any()).Return(&domain.Version{}, nil)
				variableProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("test4"))
			},
		},
		{
			name: "dataProcessService_Error",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(nil)
				versionGetService.EXPECT().Handle(ctx, gomock.Any()).Return(&domain.Version{}, nil)
				variableProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(map[string]any{}, nil)
				dataProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(nil, errors.New("test5"))
			},
			want: "test5",
		},
		{
			name: "resultRepo_Insert",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(nil)
				versionGetService.EXPECT().Handle(ctx, gomock.Any()).Return(&domain.Version{}, nil)
				variableProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(map[string]any{}, nil)
				dataProcessService.EXPECT().Handle(ctx, gomock.Any()).Return([]byte{}, nil)
				resultRepo.EXPECT().Insert(ctx, gomock.Any()).Return(int64(0), errors.New("test6"))
			},
			want: "test6",
		},
		{
			name: "taskRepo_UpdateByID_#2",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(nil)
				versionGetService.EXPECT().Handle(ctx, gomock.Any()).Return(&domain.Version{}, nil)
				variableProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(map[string]any{}, nil)
				dataProcessService.EXPECT().Handle(ctx, gomock.Any()).Return([]byte{}, nil)
				resultRepo.EXPECT().Insert(ctx, gomock.Any()).Return(int64(0), nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(errors.New("test7"))
			},
			want: "test7",
		},
		{
			name: "taskRepo_UpdateByID_#3",
			setup: func(taskRepo *MocktaskRepository, versionGetService *MockversionGetService, variableProcessService *MockvariableProcessService, dataProcessService *MockdataProcessService, resultRepo *MockresultRepository) {
				taskRepo.EXPECT().GetByID(ctx, taskID).Return(&domain.Task{}, nil)
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(nil)
				versionGetService.EXPECT().Handle(ctx, gomock.Any()).Return(&domain.Version{}, nil)
				variableProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(map[string]any{}, nil)
				dataProcessService.EXPECT().Handle(ctx, gomock.Any()).Return(nil, &domain.ProcessError{Message: "test1"})
				taskRepo.EXPECT().UpdateByID(ctx, gomock.Any()).Return(errors.New("test8"))
			},
			want: "test8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskRepo := NewMocktaskRepository(ctrl)
			versionGetService := NewMockversionGetService(ctrl)
			variableProcessService := NewMockvariableProcessService(ctrl)
			dataProcessService := NewMockdataProcessService(ctrl)
			resultRepo := NewMockresultRepository(ctrl)

			tt.setup(taskRepo, versionGetService, variableProcessService, dataProcessService, resultRepo)

			usecase := New(taskRepo, versionGetService, variableProcessService, dataProcessService, resultRepo)
			err := usecase.Handle(ctx, domain.TaskProcessIn{TaskID: taskID})
			require.ErrorContains(t, err, tt.want)
		})
	}
}
