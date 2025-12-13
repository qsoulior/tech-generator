package service

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

type Service struct {
	templateRepo   templateRepository
	versionRepo    versionRepository
	variableRepo   variableRepository
	constraintRepo constraintRepository
	trManager      trm.Manager
}

func New(
	templateRepo templateRepository,
	versionRepo versionRepository,
	variableRepo variableRepository,
	constraintRepo constraintRepository,
	trManager trm.Manager,
) *Service {
	return &Service{
		templateRepo:   templateRepo,
		versionRepo:    versionRepo,
		variableRepo:   variableRepo,
		constraintRepo: constraintRepo,
		trManager:      trManager,
	}
}

func (u *Service) Handle(ctx context.Context, in domain.VersionCreateIn) (int64, error) {
	// validate input
	if err := in.Validate(); err != nil {
		return 0, err
	}

	// create version
	var versionID int64
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		var err error
		versionID, err = u.createVersion(ctx, in)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return versionID, nil
}

func (u *Service) createVersion(ctx context.Context, in domain.VersionCreateIn) (int64, error) {
	// create version
	version := domain.Version{
		TemplateID: in.TemplateID,
		AuthorID:   in.AuthorID,
		Data:       in.Data,
	}

	versionID, err := u.versionRepo.Create(ctx, version)
	if err != nil {
		return 0, fmt.Errorf("version repo - create: %w", err)
	}

	// create variables
	err = u.createVariables(ctx, versionID, in.Variables)
	if err != nil {
		return 0, err
	}

	// update template
	templateToUpdate := domain.TemplateToUpdate{
		ID:            in.TemplateID,
		LastVersionID: versionID,
	}

	err = u.templateRepo.UpdateByID(ctx, templateToUpdate)
	if err != nil {
		return 0, fmt.Errorf("template repo - update by id: %w", err)
	}

	return versionID, nil
}

func (u *Service) createVariables(ctx context.Context, templateVersionID int64, variables []domain.Variable) error {
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
			IsInput:    v.IsInput,
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
	constraints := lo.FlatMap(variables, func(v domain.Variable, i int) []domain.ConstraintToCreate {
		return lo.Map(v.Constraints, func(c domain.Constraint, _ int) domain.ConstraintToCreate {
			return domain.ConstraintToCreate{
				VariableID: variableIDs[i],
				Name:       c.Name,
				Expression: c.Expression,
				IsActive:   c.IsActive,
			}
		})
	})

	if len(constraints) == 0 {
		return nil
	}

	err = u.constraintRepo.Create(ctx, constraints)
	if err != nil {
		return fmt.Errorf("constraint repo - create: %w", err)
	}

	return nil
}
