package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list/domain"
)

type Usecase struct {
	projectRepo     projectRepository
	projectUserRepo projectUserRepository
}

func New(projectRepo projectRepository, projectUserRepo projectUserRepository) *Usecase {
	return &Usecase{
		projectRepo:     projectRepo,
		projectUserRepo: projectUserRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectUserListIn) ([]domain.ProjectUser, error) {
	project, err := u.projectRepo.GetByID(ctx, in.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project repo - get by id: %w", err)
	}

	if project == nil {
		return nil, domain.ErrProjectNotFound
	}

	if project.AuthorID != in.UserID {
		return nil, domain.ErrProjectInvalid
	}

	users, err := u.projectUserRepo.GetByProjectID(ctx, in.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project user repo - get by project id: %w", err)
	}

	return users, nil
}
