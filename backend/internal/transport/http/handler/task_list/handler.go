package task_list_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) TaskList(ctx context.Context, params api.TaskListParams) (api.TaskListRes, error) {
	out, err := h.usecase.Handle(ctx, convertRequestToIn(params))
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("task list usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertRequestToIn(params api.TaskListParams) domain.TaskListIn {
	filter := domain.TaskListFilter{
		UserID:    params.XUserID,
		VersionID: params.VersionID,
	}

	if creatorID, ok := params.CreatorID.Get(); ok {
		filter.CreatorID = &creatorID
	}

	in := domain.TaskListIn{
		Page:   params.Page,
		Size:   params.Size,
		Filter: filter,
	}

	if sorting, ok := params.Sorting.Get(); ok {
		in.Sorting = &sorting_domain.Sorting{
			Attribute: sorting.Attribute,
			Direction: sorting_domain.SortingDirection(sorting.Direction),
		}
	}

	return in
}

func convertOutToResponse(out domain.TaskListOut) api.TaskListResponse {
	return api.TaskListResponse{
		Tasks:      lo.Map(out.Tasks, func(t domain.Task, _ int) api.TaskListResponseTasksItem { return convertTaskToResponse(t) }),
		TotalTasks: out.TotalTasks,
		TotalPages: out.TotalPages,
	}
}

func convertTaskToResponse(task domain.Task) api.TaskListResponseTasksItem {
	taskResponse := api.TaskListResponseTasksItem{
		ID:          task.ID,
		Status:      api.TaskStatus(task.Status),
		CreatorName: task.CreatorName,
		CreatedAt:   task.CreatedAt,
	}

	if task.UpdatedAt != nil {
		taskResponse.UpdatedAt.SetTo(*task.UpdatedAt)
	}

	return taskResponse
}
