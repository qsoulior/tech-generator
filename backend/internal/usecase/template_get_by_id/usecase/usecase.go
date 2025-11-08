package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

type Usecase struct {
	templateRepo           templateRepository
	templateVersionRepo    templateVersionRepository
	variableRepo           variableRepository
	variableConstraintRepo variableConstraintRepository
}

func New(
	templateRepo templateRepository,
	templateVersionRepo templateVersionRepository,
	variableRepo variableRepository,
	variableConstraintRepo variableConstraintRepository,
) *Usecase {
	return &Usecase{
		templateRepo:           templateRepo,
		templateVersionRepo:    templateVersionRepo,
		variableRepo:           variableRepo,
		variableConstraintRepo: variableConstraintRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateGetByIDIn) (*domain.TemplateGetByIDOut, error) {
	// get template
	template, err := u.getTemplate(ctx, in)
	if err != nil {
		return nil, err
	}

	if template.LastVersionID == nil {
		return &domain.TemplateGetByIDOut{VersionID: 0, VersionNumber: 0}, nil
	}

	// get last version
	version, err := u.templateVersionRepo.GetByID(ctx, *template.LastVersionID)
	if err != nil {
		return nil, fmt.Errorf("template version repo - get by id: %w", err)
	}

	if version == nil {
		return nil, domain.ErrTemplateVersionNotFound
	}

	// get variables
	variables, err := u.variableRepo.ListByVersionID(ctx, version.ID)
	if err != nil {
		return nil, fmt.Errorf("variable repo - list by version id: %w", err)
	}

	err = u.fillVariableConstraints(ctx, variables)
	if err != nil {
		return nil, err
	}

	out := domain.TemplateGetByIDOut{
		VersionID:     version.ID,
		VersionNumber: version.Number,
		CreatedAt:     version.CreatedAt,
		Data:          version.Data,
		Variables:     variables,
	}

	return &out, nil
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

func (u *Usecase) fillVariableConstraints(ctx context.Context, variables []domain.Variable) error {
	variableIDs := lo.Map(variables, func(v domain.Variable, _ int) int64 { return v.ID })

	if len(variableIDs) == 0 {
		return nil
	}

	constraints, err := u.variableConstraintRepo.ListByVariableIDs(ctx, variableIDs)
	if err != nil {
		return fmt.Errorf("variable constraint repo - list by variables ids: %w", err)
	}

	constraintsByVariable := lo.GroupBy(constraints, func(c domain.VariableConstraint) int64 { return c.VariableID })

	for i := range variables {
		variables[i].Constraints = constraintsByVariable[variables[i].ID]
	}

	return nil
}
