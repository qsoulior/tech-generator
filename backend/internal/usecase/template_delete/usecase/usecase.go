package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete/domain"
)

type Usecase struct {
	templateRepo templateRepository
}

func New(templateRepo templateRepository) *Usecase {
	return &Usecase{
		templateRepo: templateRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateDeleteIn) error {
	// get template
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return domain.ErrTemplateNotFound
	}

	if template.ProjectAuthorID != in.UserID && template.AuthorID != in.UserID {
		return domain.ErrTemplateInvalid
	}

	// delete template
	err = u.templateRepo.DeleteByID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - delete by id: %w", err)
	}

	return nil
}
