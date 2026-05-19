package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type Usecase struct {
	templateRepo templateRepository
	taskRepo     taskRepository
}

func New(templateRepo templateRepository, taskRepo taskRepository) *Usecase {
	return &Usecase{
		templateRepo: templateRepo,
		taskRepo:     taskRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TaskListIn) (*domain.TaskListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// check template access
	err := u.handleTemplate(ctx, in)
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

func (u *Usecase) handleTemplate(ctx context.Context, in domain.TaskListIn) error {
	// get template
	template, err := u.templateRepo.GetByID(ctx, in.Filter.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return domain.ErrTemplateNotFound
	}

	// check permission
	isWriter := lo.SomeBy(template.TemplateUsers, func(user domain.TemplateUser) bool {
		return user.ID == in.Filter.UserID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if template.ProjectAuthorID != in.Filter.UserID && template.TemplateAuthorID != in.Filter.UserID && !isWriter {
		return domain.ErrTemplateInvalid
	}

	return nil
}
