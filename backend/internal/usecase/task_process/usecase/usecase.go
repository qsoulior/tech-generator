package usecase

import (
	"context"
	"errors"
	"fmt"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type Usecase struct {
	taskRepo               taskRepository
	versionGetService      versionGetService
	variableProcessService variableProcessService
	dataProcessService     dataProcessService
}

func New(
	taskRepo taskRepository,
	versionGetService versionGetService,
	variableProcessService variableProcessService,
	dataProcessService dataProcessService,
) *Usecase {
	return &Usecase{
		taskRepo:               taskRepo,
		versionGetService:      versionGetService,
		variableProcessService: variableProcessService,
		dataProcessService:     dataProcessService,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TaskProcessIn) error {
	// get task
	task, err := u.taskRepo.GetByID(ctx, in.TaskID)
	if err != nil {
		return fmt.Errorf("task repo - get by id: %w", err)
	}

	if task == nil {
		return domain.ErrTaskNotFound
	}

	// update task
	taskUpdate := domain.TaskUpdate{ID: in.TaskID, Status: task_domain.StatusInProgress}
	err = u.taskRepo.UpdateByID(ctx, taskUpdate)
	if err != nil {
		return fmt.Errorf("task repo - update by id: %w", err)
	}

	// handle task
	result, err := u.handleTask(ctx, *task)
	if err != nil {
		var processErr *domain.ProcessError
		if errors.As(err, &processErr) {
			// update task
			taskUpdate = domain.TaskUpdate{ID: in.TaskID, Status: task_domain.StatusFailed, Error: processErr}
			err = u.taskRepo.UpdateByID(ctx, taskUpdate)
			if err != nil {
				return fmt.Errorf("task repo - update by id: %w", err)
			}
			return nil
		}
		return err
	}

	// update task
	taskUpdate = domain.TaskUpdate{ID: in.TaskID, Status: task_domain.StatusSucceed, Result: result}
	err = u.taskRepo.UpdateByID(ctx, taskUpdate)
	if err != nil {
		return fmt.Errorf("task repo - update by id: %w", err)
	}

	return nil
}

func (u *Usecase) handleTask(ctx context.Context, task domain.Task) ([]byte, error) {
	// get version
	version, err := u.versionGetService.Handle(ctx, task.VersionID)
	if err != nil {
		return nil, err
	}

	// process variables
	variableProcessIn := domain.VariableProcessIn{
		Variables: version.Variables,
		Payload:   task.Payload,
	}
	variableValues, err := u.variableProcessService.Handle(ctx, variableProcessIn)
	if err != nil {
		return nil, err
	}

	// process data
	dataProcessIn := domain.DataProcessIn{
		Values: variableValues,
		Data:   version.Data,
	}
	result, err := u.dataProcessService.Handle(ctx, dataProcessIn)
	if err != nil {
		return nil, err
	}

	return result, nil
}
