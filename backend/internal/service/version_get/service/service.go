package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type Service struct {
	versionRepo    versionRepository
	variableRepo   variableRepository
	constraintRepo constraintRepository
}

func New(
	versionRepo versionRepository,
	variableRepo variableRepository,
	constraintRepo constraintRepository,
) *Service {
	return &Service{
		versionRepo:    versionRepo,
		variableRepo:   variableRepo,
		constraintRepo: constraintRepo,
	}
}

func (u *Service) Handle(ctx context.Context, versionID int64) (*domain.Version, error) {
	// get version
	version, err := u.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("version repo - get by id: %w", err)
	}

	if version == nil {
		return nil, domain.ErrVersionNotFound
	}

	// get variables
	version.Variables, err = u.variableRepo.ListByVersionID(ctx, version.ID)
	if err != nil {
		return nil, fmt.Errorf("variable repo - list by version id: %w", err)
	}

	err = u.fillVariableConstraints(ctx, version.Variables)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func (u *Service) fillVariableConstraints(ctx context.Context, variables []domain.Variable) error {
	variableIDs := lo.Map(variables, func(v domain.Variable, _ int) int64 { return v.ID })

	if len(variableIDs) == 0 {
		return nil
	}

	constraints, err := u.constraintRepo.ListByVariableIDs(ctx, variableIDs)
	if err != nil {
		return fmt.Errorf("constraint repo - list by variable ids: %w", err)
	}

	constraintsByVariable := lo.GroupBy(constraints, func(c domain.Constraint) int64 { return c.VariableID })

	for i := range variables {
		variables[i].Constraints = constraintsByVariable[variables[i].ID]
	}

	return nil
}
