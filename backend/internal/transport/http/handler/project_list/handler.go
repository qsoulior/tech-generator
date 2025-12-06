package project_list_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ProjectList(ctx context.Context, params api.ProjectListParams) (api.ProjectListRes, error) {
	out, err := h.usecase.Handle(ctx, convertRequestToIn(params))
	if err != nil {
		var validationErr *error_domain.ValidationError
		if errors.As(err, &validationErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("project list by user usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertRequestToIn(params api.ProjectListParams) domain.ProjectListByUserIn {
	filter := domain.ProjectListByUserFilter{
		UserID: params.XUserID,
	}

	if projectName, ok := params.ProjectName.Get(); ok {
		filter.ProjectName = &projectName
	}

	in := domain.ProjectListByUserIn{
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

func convertOutToResponse(out domain.ProjectListByUserOut) api.ProjectListResponse {
	return api.ProjectListResponse{
		Projects: lo.Map(out.Projects, func(p domain.Project, _ int) api.ProjectListResponseProjectsItem {
			return api.ProjectListResponseProjectsItem{
				ID:         p.ID,
				Name:       p.Name,
				AuthorName: p.AuthorName,
			}
		}),
		TotalProjects: out.TotalProjects,
		TotalPages:    out.TotalPages,
	}
}
