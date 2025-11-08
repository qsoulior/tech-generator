package variable_repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_Create() {
	ctx := context.Background()
	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = true
		t.ProjectID = nil
		t.AuthorID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template version
	templateVersion := test_db.GenerateEntity(func(v *test_db.TemplateVersion) {
		v.TemplateID = templateID
		v.AuthorID = nil
		v.Number = 1
	})
	templateVersionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", templateVersion)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", templateVersionID)) }()

	// variables
	wantVariables := test_db.GenerateEntities(3, func(v *test_db.Variable, i int) {
		v.VersionID = templateVersionID
	})

	variables := lo.Map(wantVariables, func(v test_db.Variable, _ int) domain.VariableToCreate {
		return domain.VariableToCreate{
			VersionID:  templateVersionID,
			Name:       v.Name,
			Type:       variable_domain.Type(v.Type),
			Expression: v.Expression,
		}
	})

	gotIDs, err := repo.Create(ctx, variables)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "variable", gotIDs)) }()

	require.Len(s.T(), gotIDs, len(variables))

	gotVariables, err := test_db.SelectEntitiesByID[test_db.Variable](s.C(), "variable", gotIDs)
	require.NoError(s.T(), err)

	for i := range wantVariables {
		wantVariables[i].ID = gotIDs[i]
	}

	require.Equal(s.T(), wantVariables, gotVariables)
}
