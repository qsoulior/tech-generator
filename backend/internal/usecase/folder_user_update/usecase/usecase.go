package usecase

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_user_update/domain"
)

type Usecase struct {
	folderRepo     folderRepository
	folderUserRepo folderUserRepository
	userRepo       userRepository
	trManager      trm.Manager
}

func New(folderRepo folderRepository, folderUserRepo folderUserRepository, userRepo userRepository, trManager trm.Manager) *Usecase {
	return &Usecase{
		folderRepo:     folderRepo,
		folderUserRepo: folderUserRepo,
		userRepo:       userRepo,
		trManager:      trManager,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.FolderUserUpdateIn) error {
	// get folder
	folder, err := u.folderRepo.GetByID(ctx, in.FolderID)
	if err != nil {
		return fmt.Errorf("folder repo - get by id: %w", err)
	}

	if folder == nil {
		return domain.ErrFolderNotFound
	}

	if folder.RootAuthorID != in.UserID && folder.AuthorID != in.UserID {
		return domain.ErrFolderInvalid
	}

	// get existing user ids
	userIDs := lo.Map(in.Users, func(u domain.FolderUser, _ int) int64 { return u.ID })

	userExistingIDs, err := u.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("user repo - get by ids: %w", err)
	}

	userExistingIDsSet := lo.Keyify(userExistingIDs)

	folderUsersIn := lo.Filter(in.Users, func(u domain.FolderUser, _ int) bool {
		_, exists := userExistingIDsSet[u.ID]
		return exists
	})

	// get existing folder users
	folderUsersExisting, err := u.folderUserRepo.GetByFolderID(ctx, in.FolderID)
	if err != nil {
		return fmt.Errorf("folder user repo - get by folder id: %w", err)
	}

	// update folder users
	folderUsersNew := lo.Without(folderUsersIn, folderUsersExisting...)

	folderUserExistingIDs := lo.Map(folderUsersExisting, func(u domain.FolderUser, _ int) int64 { return u.ID })

	folderUserInIDs := lo.Map(folderUsersIn, func(u domain.FolderUser, _ int) int64 { return u.ID })

	folderUserMissingIDs := lo.Without(folderUserExistingIDs, folderUserInIDs...)

	err = u.trManager.Do(ctx, func(ctx context.Context) error {
		return u.updateUsers(ctx, in.FolderID, folderUsersNew, folderUserMissingIDs)
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) updateUsers(ctx context.Context, folderID int64, folderUsersNew []domain.FolderUser, folderUserMissingIDs []int64) error {
	if len(folderUsersNew) != 0 {
		err := u.folderUserRepo.Upsert(ctx, folderID, folderUsersNew)
		if err != nil {
			return fmt.Errorf("folder user repo - upsert: %w", err)
		}
	}

	if len(folderUserMissingIDs) != 0 {
		err := u.folderUserRepo.Delete(ctx, folderID, folderUserMissingIDs)
		if err != nil {
			return fmt.Errorf("folder user repo - delete: %w", err)
		}
	}

	return nil
}
