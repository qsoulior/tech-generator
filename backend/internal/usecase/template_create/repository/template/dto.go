package template_repository

import "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create/domain"

func toValues(t domain.Template) []any {
	return []any{
		t.Name,
		t.IsDefault,
		t.FolderID,
		t.AuthorID,
		t.RootAuthorID,
	}
}
