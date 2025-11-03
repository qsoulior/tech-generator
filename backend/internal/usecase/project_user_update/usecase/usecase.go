package usecase

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update/domain"
)

type Usecase struct {
	projectRepo     projectRepository
	projectUserRepo projectUserRepository
	userRepo        userRepository
	trManager       trm.Manager
}

func New(projectRepo projectRepository, projectUserRepo projectUserRepository, userRepo userRepository, trManager trm.Manager) *Usecase {
	return &Usecase{
		projectRepo:     projectRepo,
		projectUserRepo: projectUserRepo,
		userRepo:        userRepo,
		trManager:       trManager,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.ProjectUserUpdateIn) error {
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

	// get existing user ids
	userIDs := lo.Map(in.Users, func(u domain.ProjectUser, _ int) int64 { return u.ID })

	userExistingIDs, err := u.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("user repo - get by ids: %w", err)
	}

	userExistingIDsSet := lo.Keyify(userExistingIDs)

	projectUsersIn := lo.Filter(in.Users, func(u domain.ProjectUser, _ int) bool {
		_, exists := userExistingIDsSet[u.ID]
		return exists
	})

	// get existing project users
	projectUsersExisting, err := u.projectUserRepo.GetByProjectID(ctx, in.ProjectID)
	if err != nil {
		return fmt.Errorf("project user repo - get by project id: %w", err)
	}

	// update project users
	projectUsersNew := lo.Without(projectUsersIn, projectUsersExisting...)

	projectUserExistingIDs := lo.Map(projectUsersExisting, func(u domain.ProjectUser, _ int) int64 { return u.ID })

	projectUserInIDs := lo.Map(projectUsersIn, func(u domain.ProjectUser, _ int) int64 { return u.ID })

	projectUserMissingIDs := lo.Without(projectUserExistingIDs, projectUserInIDs...)

	err = u.trManager.Do(ctx, func(ctx context.Context) error {
		return u.updateUsers(ctx, in.ProjectID, projectUsersNew, projectUserMissingIDs)
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) updateUsers(ctx context.Context, projectID int64, projectUsersNew []domain.ProjectUser, projectUserMissingIDs []int64) error {
	if len(projectUsersNew) != 0 {
		err := u.projectUserRepo.Upsert(ctx, projectID, projectUsersNew)
		if err != nil {
			return fmt.Errorf("project user repo - upsert: %w", err)
		}
	}

	if len(projectUserMissingIDs) != 0 {
		err := u.projectUserRepo.Delete(ctx, projectID, projectUserMissingIDs)
		if err != nil {
			return fmt.Errorf("project user repo - delete: %w", err)
		}
	}

	return nil
}
