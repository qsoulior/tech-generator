package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	version_create_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from/domain"
)

type Usecase struct {
	templateRepo         templateRepository
	versionGetService    versionGetService
	versionCreateService versionCreateService
}

func New(templateRepo templateRepository, versionGetService versionGetService, versionCreateService versionCreateService) *Usecase {
	return &Usecase{
		templateRepo:         templateRepo,
		versionGetService:    versionGetService,
		versionCreateService: versionCreateService,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.VersionCreateFromIn) error {
	// validate template
	err := u.validateTemplate(ctx, in)
	if err != nil {
		return err
	}

	// get version
	version, err := u.versionGetService.Handle(ctx, in.VersionID)
	if err != nil {
		return err
	}

	if version.TemplateID != in.TemplateID {
		return domain.ErrVersionInvalid
	}

	// create version
	versionCreateIn := version_create_domain.VersionCreateIn{
		AuthorID:   in.AuthorID,
		TemplateID: in.TemplateID,
		Data:       version.Data,
		Variables:  convertVariables(version.Variables),
	}

	err = u.versionCreateService.Handle(ctx, versionCreateIn)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) validateTemplate(ctx context.Context, in domain.VersionCreateFromIn) error {
	// get template by id
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return domain.ErrTemplateNotFound
	}

	// check permission
	if template.ProjectAuthorID != in.AuthorID && template.AuthorID != in.AuthorID {
		return domain.ErrTemplateInvalid
	}

	return nil
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
