package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

type Usecase struct {
	projectRepo projectRepository
}

func New(projectRepo projectRepository) *Usecase {
	return &Usecase{
		projectRepo: projectRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectCreateIn) error {
	if err := in.Validate(); err != nil {
		return err
	}

	project := domain.Project{ // nolint:staticcheck
		Name:     in.Name,
		AuthorID: in.AuthorID,
	}

	err := u.projectRepo.Create(ctx, project)
	if err != nil {
		return fmt.Errorf("project repo - create: %w", err)
	}

	return nil
}
