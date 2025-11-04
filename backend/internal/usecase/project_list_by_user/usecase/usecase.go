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

	projects, err := u.projectRepo.ListByUserID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("project repo - list by user id: %w", err)
	}

	totalProjects, err := u.projectRepo.GetTotalByUserID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("project repo - count by user id: %w", err)
	}

	out := domain.ProjectListByUserOut{
		Projects:      projects,
		TotalProjects: totalProjects,
		TotalPages:    (totalProjects + in.Size - 1) / in.Size,
	}

	return &out, nil
}
