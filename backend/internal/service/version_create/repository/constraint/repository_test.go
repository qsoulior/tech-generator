package constraint_repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

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
	templateVersion := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = nil
		v.Number = 1
	})
	templateVersionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", templateVersion)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", templateVersionID)) }()

	// variables
	variables := test_db.GenerateEntities(2, func(v *test_db.Variable, i int) {
		v.VersionID = templateVersionID
	})
	variableIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "variable", variables)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "variable", variableIDs)) }()

	// variable constraints
	wantConstraints := test_db.GenerateEntities(5, func(c *test_db.Constraint, i int) {
		c.VariableID = lo.Sample(variableIDs)
	})

	constraints := lo.Map(wantConstraints, func(c test_db.Constraint, _ int) domain.ConstraintToCreate {
		return domain.ConstraintToCreate{
			VariableID: c.VariableID,
			Name:       c.Name,
			Expression: c.Expression,
			IsActive:   c.IsActive,
		}
	})

	err = repo.Create(ctx, constraints)
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntitiesByColumn(s.C(), "variable_constraint", "variable_id", variableIDs))
	}()

	gotConstraints, err := test_db.SelectEntitiesByColumn[test_db.Constraint](s.C(), "variable_constraint", "variable_id", variableIDs)
	require.NoError(s.T(), err)

	require.Len(s.T(), gotConstraints, len(wantConstraints))

	for i := range wantConstraints {
		wantConstraints[i].ID = gotConstraints[i].ID
	}

	require.Equal(s.T(), wantConstraints, gotConstraints)
}
