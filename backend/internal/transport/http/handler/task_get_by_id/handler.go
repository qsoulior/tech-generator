package task_get_by_id_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) TaskGetByID(ctx context.Context, params api.TaskGetByIDParams) (api.TaskGetByIDRes, error) {
	in := domain.TaskGetByIDIn{
		TaskID: params.TaskID,
		UserID: params.XUserID,
	}

	out, err := h.usecase.Handle(ctx, in)
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("task get by id usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertOutToResponse(out domain.TaskGetByIDOut) api.TaskGetByIDResponse {
	return api.TaskGetByIDResponse{
		Task:   convertTaskToResponse(out.Task),
		Result: out.Result,
	}
}

func convertTaskToResponse(task domain.Task) api.TaskGetByIDResponseTask {
	taskResponse := api.TaskGetByIDResponseTask{
		ID:          task.ID,
		VersionID:   task.VersionID,
		Status:      api.TaskStatus(task.Status),
		Payload:     api.TaskGetByIDResponseTaskPayload(task.Payload),
		CreatorName: task.CreatorName,
		CreatedAt:   task.CreatedAt,
	}

	if task.Error != nil {
		taskResponse.Error.SetTo(convertTaskErrorToResponse(*task.Error))
	}

	if task.UpdatedAt != nil {
		taskResponse.UpdatedAt.SetTo(*task.UpdatedAt)
	}

	return taskResponse
}

func convertTaskErrorToResponse(taskError task_domain.ProcessError) api.TaskGetByIDResponseTaskError {
	item := api.TaskGetByIDResponseTaskError{
		VariableErrors: convertVariableErrorsToResponse(taskError.VariableErrors),
	}

	if taskError.Message != "" {
		item.Message.SetTo(taskError.Message)
	}

	return item
}

func convertVariableErrorsToResponse(variableErrors []task_domain.VariableError) []api.TaskGetByIDResponseTaskErrorVariableErrorsItem {
	return lo.Map(variableErrors, func(v task_domain.VariableError, _ int) api.TaskGetByIDResponseTaskErrorVariableErrorsItem {
		item := api.TaskGetByIDResponseTaskErrorVariableErrorsItem{
			ID:               v.ID,
			Name:             v.Name,
			ConstraintErrors: convertConstraintErrorsToResponse(v.ConstraintErrors),
		}

		if v.Message != "" {
			item.Message.SetTo(v.Message)
		}

		return item
	})
}

func convertConstraintErrorsToResponse(constraintErrors []task_domain.ConstraintError) []api.TaskGetByIDResponseTaskErrorVariableErrorsItemConstraintErrorsItem {
	return lo.Map(constraintErrors, func(c task_domain.ConstraintError, _ int) api.TaskGetByIDResponseTaskErrorVariableErrorsItemConstraintErrorsItem {
		item := api.TaskGetByIDResponseTaskErrorVariableErrorsItemConstraintErrorsItem{
			ID:   c.ID,
			Name: c.Name,
		}

		if c.Message != "" {
			item.Message.SetTo(c.Message)
		}

		return item
	})
}
