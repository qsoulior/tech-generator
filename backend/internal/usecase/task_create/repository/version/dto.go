package version_repository

import (
	"github.com/samber/lo"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

type version struct {
	ProjectAuthorID  int64  `db:"project_author_id"`
	TemplateAuthorID int64  `db:"template_author_id"`
	TemplateUserID   int64  `db:"template_user_id"`
	TemplateRole     string `db:"template_user_role"`
}

type versions []version

func (vs versions) toDomain() *domain.Version {
	if len(vs) == 0 {
		return nil
	}

	return &domain.Version{
		ProjectAuthorID:  vs[0].ProjectAuthorID,
		TemplateAuthorID: vs[0].TemplateAuthorID,
		TemplateUsers: lo.Map(vs, func(v version, _ int) domain.TemplateUser {
			return domain.TemplateUser{
				ID:   v.TemplateUserID,
				Role: user_domain.Role(v.TemplateRole),
			}
		}),
	}
}
