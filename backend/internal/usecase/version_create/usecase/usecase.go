package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create/domain"
)

type Usecase struct {
	templateRepo         templateRepository
	versionCreateService versionCreateService
}

func New(templateRepo templateRepository, versionCreateService versionCreateService) *Usecase {
	return &Usecase{
		templateRepo:         templateRepo,
		versionCreateService: versionCreateService,
	}
}

func (u *Usecase) Handle(ctx context.Context, in version_create_domain.VersionCreateIn) error {
	// get template
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return domain.ErrTemplateNotFound
	}

	// check permission
	isWriter := lo.SomeBy(template.Users, func(user domain.TemplateUser) bool {
		return user.ID == in.AuthorID && user.Role == user_domain.RoleWrite
	})

	if template.ProjectAuthorID != in.AuthorID && template.AuthorID != in.AuthorID && !isWriter {
		return domain.ErrTemplateInvalid
	}

	// create version
	err = u.versionCreateService.Handle(ctx, in)
	if err != nil {
		return err
	}

	return nil
}
