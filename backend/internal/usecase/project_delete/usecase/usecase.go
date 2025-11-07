package usecase

import (
	"context"
	"fmt"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete/domain"
)

type Usecase struct {
	projectRepo projectRepository
}

func New(projectRepo projectRepository) *Usecase {
	return &Usecase{
		projectRepo: projectRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectDeleteIn) error {
	// get project
	project, err := u.projectRepo.GetByID(ctx, in.ProjectID)
	if err != nil {
		return fmt.Errorf("project repo - get by id: %w", err)
	}

	if project == nil {
		return domain.ErrProjectNotFound
	}

	if project.AuthorID != in.UserID {
		return domain.ErrProjectInvalid
	}

	// delete template
	err = u.projectRepo.DeleteByID(ctx, in.ProjectID)
	if err != nil {
		return fmt.Errorf("project repo - delete by id: %w", err)
	}

	return nil
}
