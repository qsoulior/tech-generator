package usecase

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

type Usecase struct {
	templateRepo           templateRepository
	templateVersionRepo    templateVersionRepository
	variableRepo           variableRepository
	variableConstraintRepo variableConstraintRepository
	trManager              trm.Manager
}

func New(
	templateRepo templateRepository,
	templateVersionRepo templateVersionRepository,
	variableRepo variableRepository,
	variableConstraintRepo variableConstraintRepository,
	trManager trm.Manager,
) *Usecase {
	return &Usecase{
		templateRepo:           templateRepo,
		templateVersionRepo:    templateVersionRepo,
		variableRepo:           variableRepo,
		variableConstraintRepo: variableConstraintRepo,
		trManager:              trManager,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateVersionCreateIn) error {
	// validate input
	if err := in.Validate(); err != nil {
		return err
	}

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

	// create template version
	err = u.trManager.Do(ctx, func(ctx context.Context) error { return u.createTemplateVersion(ctx, in) })
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) createTemplateVersion(ctx context.Context, in domain.TemplateVersionCreateIn) error {
	// create template version
	templateVersion := domain.TemplateVersion{
		TemplateID: in.TemplateID,
		AuthorID:   in.AuthorID,
		Data:       in.Data,
	}

	templateVersionID, err := u.templateVersionRepo.Create(ctx, templateVersion)
	if err != nil {
		return fmt.Errorf("template version repo - create: %w", err)
	}

	// create variables
	err = u.createVariables(ctx, templateVersionID, in.Variables)
	if err != nil {
		return err
	}

	// update template
	templateToUpdate := domain.TemplateToUpdate{
		ID:            in.TemplateID,
		LastVersionID: templateVersionID,
	}

	err = u.templateRepo.UpdateByID(ctx, templateToUpdate)
	if err != nil {
		return fmt.Errorf("template repo - update by id: %w", err)
	}

	return nil
}

func (u *Usecase) createVariables(ctx context.Context, templateVersionID int64, variables []domain.Variable) error {
	if len(variables) == 0 {
		return nil
	}

	// create variables
	variablesToCreate := lo.Map(variables, func(v domain.Variable, _ int) domain.VariableToCreate {
		return domain.VariableToCreate{
			VersionID:  templateVersionID,
			Name:       v.Name,
			Type:       v.Type,
			Expression: v.Expression,
		}
	})

	variableIDs, err := u.variableRepo.Create(ctx, variablesToCreate)
	if err != nil {
		return fmt.Errorf("variable repo - create: %w", err)
	}

	if len(variables) != len(variableIDs) {
		return domain.ErrVariableIDsInvalid
	}

	// create variable constraints
	constrains := lo.FlatMap(variables, func(v domain.Variable, i int) []domain.VariableConstraintToCreate {
		return lo.Map(v.Constraints, func(c domain.VariableConstraint, _ int) domain.VariableConstraintToCreate {
			return domain.VariableConstraintToCreate{
				VariableID: variableIDs[i],
				Name:       c.Name,
				Expression: c.Expression,
				IsActive:   c.IsActive,
			}
		})
	})

	if len(constrains) == 0 {
		return nil
	}

	err = u.variableConstraintRepo.Create(ctx, constrains)
	if err != nil {
		return fmt.Errorf("variable constraint repo - create: %w", err)
	}

	return nil
}
