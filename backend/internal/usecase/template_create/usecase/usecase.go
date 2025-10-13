package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"
)

type Usecase struct {
	folderRepo   folderRepository
	templateRepo templateRepository
}

func New(folderRepo folderRepository, templateRepo templateRepository) *Usecase {
	return &Usecase{
		folderRepo:   folderRepo,
		templateRepo: templateRepo,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateCreateIn) error {
	if err := in.Validate(); err != nil {
		return err
	}

	folder, err := u.folderRepo.GetByID(ctx, in.FolderID)
	if err != nil {
		return fmt.Errorf("folder repo - get by id: %w", err)
	}

	if folder == nil {
		return domain.ErrFolderNotFound
	}

	isWriter := lo.SomeBy(folder.Users, func(user domain.FolderUser) bool {
		return user.ID == in.AuthorID && user.Role == user_domain.RoleWrite
	})

	if folder.RootAuthorID != in.AuthorID && folder.AuthorID != in.AuthorID && !isWriter {
		return domain.ErrFolderInvalid
	}

	template := domain.Template{
		Name:         in.Name,
		IsDefault:    false,
		FolderID:     in.FolderID,
		AuthorID:     in.AuthorID,
		RootAuthorID: folder.RootAuthorID,
	}

	err = u.templateRepo.Create(ctx, template)
	if err != nil {
		return fmt.Errorf("template repo - create: %w", err)
	}

	return nil
}
