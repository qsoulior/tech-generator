package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default/domain"
)

type Usecase struct {
	projectRepo          projectRepository
	sourceTemplateRepo   sourceTemplateRepository
	newTemplateRepo      newTemplateRepository
	versionGetService    versionGetService
	versionCreateService versionCreateService
}

func New(
	projectRepo projectRepository,
	sourceTemplateRepo sourceTemplateRepository,
	newTemplateRepo newTemplateRepository,
	versionGetService versionGetService,
	versionCreateService versionCreateService,
) *Usecase {
	return &Usecase{
		projectRepo:          projectRepo,
		sourceTemplateRepo:   sourceTemplateRepo,
		newTemplateRepo:      newTemplateRepo,
		versionGetService:    versionGetService,
		versionCreateService: versionCreateService,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateCreateFromDefaultIn) (*domain.TemplateCreateFromDefaultOut, error) {
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

	source, err := u.sourceTemplateRepo.GetByID(ctx, in.SourceTemplateID)
	if err != nil {
		return nil, fmt.Errorf("source template repo - get by id: %w", err)
	}

	if source == nil {
		return nil, domain.ErrSourceTemplateNotFound
	}

	if !source.IsDefault {
		return nil, domain.ErrSourceTemplateInvalid
	}

	templateID, err := u.newTemplateRepo.Create(ctx, domain.Template{
		Name:      in.Name,
		IsDefault: false,
		ProjectID: in.ProjectID,
		AuthorID:  in.AuthorID,
	})
	if err != nil {
		return nil, fmt.Errorf("new template repo - create: %w", err)
	}

	if source.LastVersionID != nil {
		version, err := u.versionGetService.Handle(ctx, *source.LastVersionID)
		if err != nil {
			return nil, err
		}

		versionIn := version_create_domain.VersionCreateIn{
			AuthorID:   in.AuthorID,
			TemplateID: templateID,
			Data:       version.Data,
			Variables:  convertVariables(version.Variables),
		}

		_, err = u.versionCreateService.Handle(ctx, versionIn)
		if err != nil {
			return nil, err
		}
	}

	return &domain.TemplateCreateFromDefaultOut{ID: templateID}, nil
}

func convertVariables(variables []version_get_domain.Variable) []version_create_domain.Variable {
	return lo.Map(variables, func(v version_get_domain.Variable, _ int) version_create_domain.Variable {
		return version_create_domain.Variable{
			Name:        v.Name,
			Type:        v.Type,
			Expression:  v.Expression,
			IsInput:     v.IsInput,
			Constraints: convertConstraints(v.Constraints),
		}
	})
}

func convertConstraints(constraints []version_get_domain.Constraint) []version_create_domain.Constraint {
	return lo.Map(constraints, func(c version_get_domain.Constraint, _ int) version_create_domain.Constraint {
		return version_create_domain.Constraint{
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})
}
