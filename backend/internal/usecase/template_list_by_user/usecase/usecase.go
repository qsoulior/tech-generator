package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user/domain"
)

type Usecase struct {
	templateRepo templateRepository
}

func New(templateRepo templateRepository) *Usecase {
	return &Usecase{
		templateRepo: templateRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateListByUserIn) (*domain.TemplateListByUserOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	templates, err := u.templateRepo.ListByUserID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("template repo - list by author id: %w", err)
	}

	totalTemplates, err := u.templateRepo.GetTotalByUserID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("template repo - count by user id: %w", err)
	}

	out := domain.TemplateListByUserOut{
		Templates:      templates,
		TotalTemplates: totalTemplates,
		TotalPages:     (totalTemplates + in.Size - 1) / in.Size,
	}

	return &out, nil
}
