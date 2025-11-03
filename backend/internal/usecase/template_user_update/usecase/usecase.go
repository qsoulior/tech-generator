package usecase

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/samber/lo"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update/domain"
)

type Usecase struct {
	templateRepo     templateRepository
	templateUserRepo templateUserRepository
	userRepo         userRepository
	trManager        trm.Manager
}

func New(templateRepo templateRepository, templateUserRepo templateUserRepository, userRepo userRepository, trManager trm.Manager) *Usecase {
	return &Usecase{
		templateRepo:     templateRepo,
		templateUserRepo: templateUserRepo,
		userRepo:         userRepo,
		trManager:        trManager,
	}
}

func (u *Usecase) Handle(ctx context.Context, in domain.TemplateUserUpdateIn) error {
	// get template
	template, err := u.templateRepo.GetByID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template repo - get by id: %w", err)
	}

	if template == nil {
		return domain.ErrTemplateNotFound
	}

	if template.ProjectAuthorID != in.UserID && template.AuthorID != in.UserID {
		return domain.ErrTemplateInvalid
	}

	// get existing user ids
	userIDs := lo.Map(in.Users, func(u domain.TemplateUser, _ int) int64 { return u.ID })

	userExistingIDs, err := u.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("user repo - get by ids: %w", err)
	}

	userExistingIDsSet := lo.Keyify(userExistingIDs)

	templateUsersIn := lo.Filter(in.Users, func(u domain.TemplateUser, _ int) bool {
		_, exists := userExistingIDsSet[u.ID]
		return exists
	})

	// get existing template users
	templateUsersExisting, err := u.templateUserRepo.GetByTemplateID(ctx, in.TemplateID)
	if err != nil {
		return fmt.Errorf("template user repo - get by template id: %w", err)
	}

	// update template users
	templateUsersNew := lo.Without(templateUsersIn, templateUsersExisting...)

	templateUserExistingIDs := lo.Map(templateUsersExisting, func(u domain.TemplateUser, _ int) int64 { return u.ID })

	templateUserInIDs := lo.Map(templateUsersIn, func(u domain.TemplateUser, _ int) int64 { return u.ID })

	templateUserMissingIDs := lo.Without(templateUserExistingIDs, templateUserInIDs...)

	err = u.trManager.Do(ctx, func(ctx context.Context) error {
		return u.updateUsers(ctx, in.TemplateID, templateUsersNew, templateUserMissingIDs)
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) updateUsers(ctx context.Context, templateID int64, templateUsersNew []domain.TemplateUser, templateUserMissingIDs []int64) error {
	if len(templateUsersNew) != 0 {
		err := u.templateUserRepo.Upsert(ctx, templateID, templateUsersNew)
		if err != nil {
			return fmt.Errorf("template user repo - upsert: %w", err)
		}
	}

	if len(templateUserMissingIDs) != 0 {
		err := u.templateUserRepo.Delete(ctx, templateID, templateUserMissingIDs)
		if err != nil {
			return fmt.Errorf("template user repo - delete: %w", err)
		}
	}

	return nil
}
