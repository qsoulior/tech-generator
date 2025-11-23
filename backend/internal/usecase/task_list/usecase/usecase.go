package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type Usecase struct {
	versionRepo versionRepository
	taskRepo    taskRepository
}

func New(versionRepo versionRepository, taskRepo taskRepository) *Usecase {
	return &Usecase{
		versionRepo: versionRepo,
		taskRepo:    taskRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TaskListIn) (*domain.TaskListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// check version
	err := u.handleVersion(ctx, in)
	if err != nil {
		return nil, err
	}

	// list tasks
	tasks, err := u.taskRepo.List(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("task repo - list: %w", err)
	}

	// get total tasks
	totalTasks, err := u.taskRepo.GetTotal(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("task repo - get total: %w", err)
	}

	out := domain.TaskListOut{
		Tasks:      tasks,
		TotalTasks: totalTasks,
		TotalPages: (totalTasks + in.Size - 1) / in.Size,
	}

	return &out, nil
}

func (u *Usecase) handleVersion(ctx context.Context, in domain.TaskListIn) error {
	// get version
	version, err := u.versionRepo.GetByID(ctx, in.Filter.VersionID)
	if err != nil {
		return fmt.Errorf("version repo - get by id: %w", err)
	}

	if version == nil {
		return domain.ErrVersionNotFound
	}

	// check permission
	isWriter := lo.SomeBy(version.TemplateUsers, func(user domain.TemplateUser) bool {
		return user.ID == in.Filter.UserID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if version.ProjectAuthorID != in.Filter.UserID && version.TemplateAuthorID != in.Filter.UserID && !isWriter {
		return domain.ErrVersionInvalid
	}

	return nil
}
