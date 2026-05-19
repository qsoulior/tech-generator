package template_repository

import (
	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_list/domain"
)

type template struct {
	ProjectAuthorID  int64   `db:"project_author_id"`
	TemplateAuthorID int64   `db:"template_author_id"`
	TemplateUserID   *int64  `db:"template_user_id"`
	TemplateRole     *string `db:"template_user_role"`
}

type templates []template

func (ts templates) toDomain() *domain.Template {
	if len(ts) == 0 {
		return nil
	}

	users := lo.FilterMap(ts, func(t template, _ int) (domain.TemplateUser, bool) {
		if t.TemplateUserID == nil {
			return domain.TemplateUser{}, false
		}
		return domain.TemplateUser{ID: *t.TemplateUserID, Role: user_domain.Role(*t.TemplateRole)}, true
	})

	return &domain.Template{
		ProjectAuthorID:  ts[0].ProjectAuthorID,
		TemplateAuthorID: ts[0].TemplateAuthorID,
		TemplateUsers:    users,
	}
}
