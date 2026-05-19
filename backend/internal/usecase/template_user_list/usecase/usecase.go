package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list/domain"
)

type Usecase struct {
	templateRepo     templateRepository
	templateUserRepo templateUserRepository
}

func New(templateRepo templateRepository, templateUserRepo templateUserRepository) *Usecase {
	return &Usecase{
		templateRepo:     templateRepo,
		templateUserRepo: templateUserRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateUserListIn) ([]domain.TemplateUser, error) {
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return nil, domain.ErrTemplateNotFound
	}

	if template.AuthorID != in.UserID && template.ProjectAuthorID != in.UserID {
		return nil, domain.ErrTemplateInvalid
	}

	users, err := u.templateUserRepo.GetByTemplateID(ctx, in.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template user repo - get by template id: %w", err)
	}

	return users, nil
}
