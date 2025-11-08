package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

type Usecase struct {
	templateRepo templateRepository
	versionRepo  versionRepository
}

func New(templateRepo templateRepository, templateVersionRepo versionRepository) *Usecase {
	return &Usecase{
		templateRepo: templateRepo,
		versionRepo:  templateVersionRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.VersionListIn) (*domain.VersionListOut, error) {
	// get template
	err := u.validateTemplate(ctx, in)
	if err != nil {
		return nil, err
	}

	versions, err := u.versionRepo.ListByTemplateID(ctx, in.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("version repo - list by template id: %w", err)
	}

	out := domain.VersionListOut{
		Versions: versions,
	}
	return &out, nil
}

func (u *Usecase) validateTemplate(ctx context.Context, in domain.VersionListIn) error {
	// get template by id
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return domain.ErrTemplateNotFound
	}

	// check permission
	if template.ProjectAuthorID != in.UserID && template.AuthorID != in.UserID {
		return domain.ErrTemplateInvalid
	}

	return nil
}
