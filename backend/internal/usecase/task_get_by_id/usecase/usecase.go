package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id/domain"
)

type Usecase struct {
	taskRepo    taskRepository
	versionRepo versionRepository
	resultRepo  resultRepository
}

func New(taskRepo taskRepository, versionRepo versionRepository, resultRepo resultRepository) *Usecase {
	return &Usecase{
		taskRepo:    taskRepo,
		versionRepo: versionRepo,
		resultRepo:  resultRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TaskGetByIDIn) (*domain.TaskGetByIDOut, error) {
	task, err := u.taskRepo.GetByID(ctx, in.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task repo - get by id: %w", err)
	}

	if task == nil {
		return nil, domain.ErrTaskNotFound
	}

	err = u.handleVersion(ctx, task.VersionID, in.UserID)
	if err != nil {
		return nil, err
	}

	if task.ResultID == nil {
		return &domain.TaskGetByIDOut{Task: *task, Result: nil}, nil
	}

	result, err := u.resultRepo.GetDataByID(ctx, *task.ResultID)
	if err != nil {
		return nil, fmt.Errorf("result repo - get data by id: %w", err)
	}

	return &domain.TaskGetByIDOut{Task: *task, Result: result}, nil
}

func (u *Usecase) handleVersion(ctx context.Context, versionID, userID int64) error {
	// get version
	version, err := u.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return fmt.Errorf("version repo - get by id: %w", err)
	}

	if version == nil {
		return domain.ErrTaskNotFound
	}

	// check permission
	isWriter := lo.SomeBy(version.TemplateUsers, func(user domain.TemplateUser) bool {
		return user.ID == userID && (user.Role == user_domain.RoleRead || user.Role == user_domain.RoleWrite)
	})

	if version.ProjectAuthorID != userID && version.TemplateAuthorID != userID && !isWriter {
		return domain.ErrTaskInvalid
	}

	return nil
}
