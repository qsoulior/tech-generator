package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user/domain"
)

type Usecase struct {
	projectRepo projectRepository
}

func New(projectRepo projectRepository) *Usecase {
	return &Usecase{
		projectRepo: projectRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectListByUserIn) (*domain.ProjectListByUserOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	projectsOwned, err := u.projectRepo.ListByAuthorID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("project repo - list by author id: %w", err)
	}

	projectsShared, err := u.projectRepo.ListByProjectUserID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("project repo - list by project user id: %w", err)
	}

	out := domain.ProjectListByUserOut{
		Owned:  projectsOwned,
		Shared: projectsShared,
	}
	return &out, nil
}
