package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_meta_by_id/domain"
)

type Usecase struct {
	templateRepo templateRepository
}

func New(templateRepo templateRepository) *Usecase {
	return &Usecase{
		templateRepo: templateRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateGetMetaByIDIn) (*domain.TemplateGetMetaByIDOut, error) {
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return nil, domain.ErrTemplateNotFound
	}

	isReader := lo.SomeBy(template.Users, func(user domain.TemplateUser) bool {
		return user.ID == in.UserID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if template.ProjectAuthorID != in.UserID && template.AuthorID != in.UserID && !isReader {
		return nil, domain.ErrTemplateInvalid
	}

	return &domain.TemplateGetMetaByIDOut{Name: template.Name}, nil
}
