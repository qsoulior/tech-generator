package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_update/domain"
)

type Usecase struct {
	projectRepo projectRepository
}

func New(projectRepo projectRepository) *Usecase {
	return &Usecase{
		projectRepo: projectRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectUpdateIn) error {
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

	isMaintainer := lo.SomeBy(project.Users, func(user domain.ProjectUser) bool {
		return user.ID == in.UserID && user.Role == user_domain.RoleMaintain
	})

	if project.AuthorID != in.UserID && !isMaintainer {
		return domain.ErrProjectInvalid
	}

	err = u.projectRepo.UpdateByID(ctx, in.ProjectID, in.Name)
	if err != nil {
		return fmt.Errorf("project repo - update by id: %w", err)
	}

	return nil
}
