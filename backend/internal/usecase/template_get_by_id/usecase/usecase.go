package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

type Usecase struct {
	templateRepo      templateRepository
	versionGetService versionGetService
}

func New(templateRepo templateRepository, versionGetService versionGetService) *Usecase {
	return &Usecase{
		templateRepo:      templateRepo,
		versionGetService: versionGetService,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateGetByIDIn) (*domain.TemplateGetByIDOut, error) {
	// get template
	template, err := u.getTemplate(ctx, in)
	if err != nil {
		return nil, err
	}

	if template.LastVersionID == nil {
		return &domain.TemplateGetByIDOut{Version: nil}, nil
	}

	// get last version
	version, err := u.versionGetService.Handle(ctx, *template.LastVersionID)
	if err != nil {
		return nil, fmt.Errorf("version get service - handle: %w", err)
	}

	return &domain.TemplateGetByIDOut{Version: version}, nil
}

func (u *Usecase) getTemplate(ctx context.Context, in domain.TemplateGetByIDIn) (*domain.Template, error) {
	// get template by id
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return nil, domain.ErrTemplateNotFound
	}

	// check permission
	isReader := lo.SomeBy(template.Users, func(user domain.TemplateUser) bool {
		return user.ID == in.UserID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if template.ProjectAuthorID != in.UserID && template.AuthorID != in.UserID && !isReader {
		return nil, domain.ErrTemplateInvalid
	}

	return template, nil
}
