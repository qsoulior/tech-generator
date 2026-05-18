package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

type Usecase struct {
	templateRepo templateRepository
}

func New(templateRepo templateRepository) *Usecase {
	return &Usecase{
		templateRepo: templateRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateListDefaultIn) (*domain.TemplateListDefaultOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	templates, err := u.templateRepo.ListDefault(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("template repo - list default: %w", err)
	}

	totalTemplates, err := u.templateRepo.GetTotalDefault(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("template repo - count default: %w", err)
	}

	out := domain.TemplateListDefaultOut{
		Templates:      templates,
		TotalTemplates: totalTemplates,
		TotalPages:     (totalTemplates + in.Size - 1) / in.Size,
	}

	return &out, nil
}
