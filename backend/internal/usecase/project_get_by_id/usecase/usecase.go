package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id/domain"
)

type Usecase struct {
	projectRepo projectRepository
}

func New(projectRepo projectRepository) *Usecase {
	return &Usecase{
		projectRepo: projectRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectGetByIDIn) (*domain.ProjectGetByIDOut, error) {
	project, err := u.projectRepo.GetByID(ctx, in.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project repo - get by id: %w", err)
	}

	if project == nil {
		return nil, domain.ErrProjectNotFound
	}

	isReader := lo.SomeBy(project.Users, func(user domain.ProjectUser) bool {
		return user.ID == in.UserID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if project.AuthorID != in.UserID && !isReader {
		return nil, domain.ErrProjectInvalid
	}

	return &domain.ProjectGetByIDOut{
		Name:       project.Name,
		AuthorName: project.AuthorName,
	}, nil
}
