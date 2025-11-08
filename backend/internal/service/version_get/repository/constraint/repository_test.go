package constraint_repository

import (
	"context"
	"slices"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_ListByVariableIDs() {
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
	version := test_db.GenerateEntity(func(v *test_db.TemplateVersion) {
		v.TemplateID = templateID
		v.AuthorID = nil
	})
	versionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", version)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", versionID)) }()

	// variables
	variables := test_db.GenerateEntities(3, func(v *test_db.Variable, i int) {
		v.VersionID = versionID
	})
	variableIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "variable", variables)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "variable", variableIDs)) }()

	// variable constraints
	constraints := slices.Concat(
		test_db.GenerateEntities(3, func(c *test_db.VariableConstraint, i int) {
			c.VariableID = variableIDs[0]
		}),
		test_db.GenerateEntities(2, func(c *test_db.VariableConstraint, i int) {
			c.VariableID = variableIDs[1]
		}),
		test_db.GenerateEntities(2, func(c *test_db.VariableConstraint, i int) {
			c.VariableID = variableIDs[2]
		}),
	)
	constraintIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "variable_constraint", constraints)
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "variable_constraint", constraintIDs))
	}()

	got, err := repo.ListByVariableIDs(ctx, variableIDs[:2])
	require.NoError(s.T(), err)

	want := lo.Map(constraints[:5], func(c test_db.VariableConstraint, _ int) domain.Constraint {
		return domain.Constraint{
			ID:         c.ID,
			VariableID: c.VariableID,
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})
	slices.SortFunc(want, func(a, b domain.Constraint) int { return int(a.ID - b.ID) })
	require.Equal(s.T(), want, got)
}
