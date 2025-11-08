package variable_repository

import (
	"context"
	"slices"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_ListByVersionID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.ProjectID = nil
		t.AuthorID = nil
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template version
	versions := test_db.GenerateEntities(2, func(v *test_db.TemplateVersion, _ int) {
		v.TemplateID = templateID
		v.AuthorID = nil
	})
	versionIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template_version", versions)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template_version", versionIDs)) }()

	// variables
	variables := slices.Concat(
		test_db.GenerateEntities(3, func(v *test_db.Variable, i int) {
			v.VersionID = versionIDs[0]
		}),
		test_db.GenerateEntities(2, func(v *test_db.Variable, i int) {
			v.VersionID = versionIDs[1]
		}),
	)
	variableIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "variable", variables)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "variable", variableIDs)) }()

	got, err := repo.ListByVersionID(ctx, versionIDs[0])
	require.NoError(s.T(), err)

	want := lo.Map(variables[:3], func(v test_db.Variable, _ int) domain.Variable {
		return domain.Variable{
			ID:         v.ID,
			Name:       v.Name,
			Type:       variable_domain.Type(v.Type),
			Expression: v.Expression,
		}
	})
	slices.SortFunc(want, func(a, b domain.Variable) int { return int(a.ID - b.ID) })
	require.Equal(s.T(), want, got)
}
