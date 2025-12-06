package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

type Usecase struct {
	versionRepo versionRepository
	taskRepo    taskRepository
	publisher   publisher
}

func New(versionRepo versionRepository, taskRepo taskRepository, publisher publisher) *Usecase {
	return &Usecase{
		versionRepo: versionRepo,
		taskRepo:    taskRepo,
		publisher:   publisher,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TaskCreateIn) error {
	// check version
	err := u.handleVersion(ctx, in)
	if err != nil {
		return err
	}

	// insert task
	taskID, err := u.taskRepo.Insert(ctx, in)
	if err != nil {
		return fmt.Errorf("task repo - insert: %w", err)
	}

	// send message
	err = u.publisher.PublishTaskCreated(ctx, taskID)
	if err != nil {
		return fmt.Errorf("publisher - publish task created: %w", err)
	}

	return nil
}

func (u *Usecase) handleVersion(ctx context.Context, in domain.TaskCreateIn) error {
	// get version
	version, err := u.versionRepo.GetByID(ctx, in.VersionID)
	if err != nil {
		return fmt.Errorf("version repo - get by id: %w", err)
	}

	if version == nil {
		return domain.ErrVersionNotFound
	}

	// check permission
	isWriter := lo.SomeBy(version.TemplateUsers, func(user domain.TemplateUser) bool {
		return user.ID == in.CreatorID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if version.ProjectAuthorID != in.CreatorID && version.TemplateAuthorID != in.CreatorID && !isWriter {
		return domain.ErrVersionInvalid
	}

	return nil
}
