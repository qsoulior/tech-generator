package version_list_handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

type Handler struct {
	usecase usecase
}

func (h *Handler) VersionList(ctx context.Context, params api.VersionListParams) (api.VersionListRes, error) {
	out, err := h.usecase.Handle(ctx, convertRequestToIn(params))
	if err != nil {
		var baseErr *error_domain.BaseError
		if errors.As(err, &baseErr) {
			return &api.Error{Message: err.Error()}, nil
		}

		return nil, fmt.Errorf("version list usecase: %w", err)
	}

	resp := convertOutToResponse(*out)
	return &resp, nil
}

func convertRequestToIn(params api.VersionListParams) domain.VersionListIn {
	return domain.VersionListIn{
		TemplateID: params.TemplateID,
		UserID:     params.XUserID,
	}
}

func convertOutToResponse(out domain.VersionListOut) api.VersionListResponse {
	return api.VersionListResponse{
		Versions: lo.Map(out.Versions, func(v domain.Version, _ int) api.VersionListResponseVersionsItem {
			return api.VersionListResponseVersionsItem{
				ID:         v.ID,
				Number:     v.Number,
				AuthorName: v.AuthorName,
				CreatedAt:  v.CreatedAt,
			}
		}),
	}
}
