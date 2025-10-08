package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/folder_create/domain"
)

type Usecase struct {
	folderRepo folderRepository
}

func New(folderRepo folderRepository) *Usecase {
	return &Usecase{
		folderRepo: folderRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.FolderCreateIn) error {
	if err := in.Validate(); err != nil {
		return err
	}

	if in.ParentID == nil {
		return u.createWithoutParent(ctx, in)
	}

	return u.createWithParent(ctx, in)
}

func (u *Usecase) createWithoutParent(ctx context.Context, in domain.FolderCreateIn) error {
	err := u.folderRepo.Create(ctx, in.Name, in.AuthorID, in.AuthorID, in.ParentID)
	if err != nil {
		return fmt.Errorf("folder repo - create: %w", err)
	}

	return nil
}

func (u *Usecase) createWithParent(ctx context.Context, in domain.FolderCreateIn) error {
	folder, err := u.folderRepo.GetByID(ctx, *in.ParentID)
	if err != nil {
		return fmt.Errorf("folder repo - get by id: %w", err)
	}

	if folder == nil {
		return domain.ErrParentNotFound
	}

	isMaintainer := lo.SomeBy(folder.Users, func(user domain.FolderUser) bool {
		return user.ID == in.AuthorID && user.Role == user_domain.RoleMaintain
	})

	if folder.RootAuthorID != in.AuthorID && folder.AuthorID != in.AuthorID && !isMaintainer {
		return domain.ErrParentInvalid
	}

	err = u.folderRepo.Create(ctx, in.Name, in.AuthorID, folder.RootAuthorID, in.ParentID)
	if err != nil {
		return fmt.Errorf("folder repo - create: %w", err)
	}

	return nil
}
