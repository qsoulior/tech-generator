package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_import/domain"
)

type Usecase struct {
	projectRepo          projectRepository
	templateRepo         templateRepository
	versionCreateService versionCreateService
}

func New(projectRepo projectRepository, templateRepo templateRepository, versionCreateService versionCreateService) *Usecase {
	return &Usecase{
		projectRepo:          projectRepo,
		templateRepo:         templateRepo,
		versionCreateService: versionCreateService,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateImportIn) (*domain.TemplateImportOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	project, err := u.projectRepo.GetByID(ctx, in.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project repo - get by id: %w", err)
	}

	if project == nil {
		return nil, domain.ErrProjectNotFound
	}

	isWriter := lo.SomeBy(project.Users, func(user domain.ProjectUser) bool {
		return user.ID == in.AuthorID && user.Role == user_domain.RoleWrite
	})

	if project.AuthorID != in.AuthorID && !isWriter {
		return nil, domain.ErrProjectInvalid
	}

	templateID, err := u.templateRepo.Create(ctx, domain.Template{
		Name:      in.Name,
		IsDefault: false,
		ProjectID: in.ProjectID,
		AuthorID:  in.AuthorID,
	})
	if err != nil {
		return nil, fmt.Errorf("template repo - create: %w", err)
	}

	if in.Version != nil {
		versionIn := version_create_domain.VersionCreateIn{
			AuthorID:   in.AuthorID,
			TemplateID: templateID,
			Data:       in.Version.Data,
			Variables:  convertVariables(in.Version.Variables),
		}

		_, err = u.versionCreateService.Handle(ctx, versionIn)
		if err != nil {
			return nil, err
		}
	}

	return &domain.TemplateImportOut{ID: templateID}, nil
}

func convertVariables(variables []domain.Variable) []version_create_domain.Variable {
	return lo.Map(variables, func(v domain.Variable, _ int) version_create_domain.Variable {
		return version_create_domain.Variable{
			Name:        v.Name,
			Title:       v.Title,
			Type:        v.Type,
			Expression:  v.Expression,
			IsInput:     v.IsInput,
			Constraints: convertConstraints(v.Constraints),
		}
	})
}

func convertConstraints(constraints []domain.Constraint) []version_create_domain.Constraint {
	return lo.Map(constraints, func(c domain.Constraint, _ int) version_create_domain.Constraint {
		return version_create_domain.Constraint{
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})
}
