package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

type Usecase struct {
	projectRepo  projectRepository
	templateRepo templateRepository
}

func New(projectRepo projectRepository, templateRepo templateRepository) *Usecase {
	return &Usecase{
		projectRepo:  projectRepo,
		templateRepo: templateRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateCreateIn) error {
	if err := in.Validate(); err != nil {
		return err
	}

	project, err := u.projectRepo.GetByID(ctx, in.ProjectID)
	if err != nil {
		return fmt.Errorf("project repo - get by id: %w", err)
	}

	if project == nil {
		return domain.ErrProjectNotFound
	}

	isWriter := lo.SomeBy(project.Users, func(user domain.ProjectUser) bool {
		return user.ID == in.AuthorID && user.Role == user_domain.RoleWrite
	})

	if project.AuthorID != in.AuthorID && !isWriter {
		return domain.ErrProjectInvalid
	}

	template := domain.Template{
		Name:      in.Name,
		IsDefault: false,
		ProjectID: in.ProjectID,
		AuthorID:  in.AuthorID,
	}

	err = u.templateRepo.Create(ctx, template)
	if err != nil {
		return fmt.Errorf("template repo - create: %w", err)
	}

	return nil
}
