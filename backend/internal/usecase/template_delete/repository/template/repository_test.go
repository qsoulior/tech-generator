package template_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_GetByID() {
	ctx := context.Background()

	repo := New(s.C().DB())

	s.T().Run("Exists", func(t *testing.T) {
		// users
		users := test_db.GenerateEntities[test_db.User](4)
		userIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "usr", users)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntitiesByID(s.C(), "usr", userIDs)) }()

		// project
		project := test_db.GenerateEntity(func(p *test_db.Project) {
			p.AuthorID = users[0].ID
		})
		projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

		// template
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.IsDefault = false
			t.ProjectID = &projectID
			t.AuthorID = &users[1].ID
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)

		want := domain.Template{
			AuthorID:        *template.AuthorID,
			ProjectAuthorID: project.AuthorID,
		}
		require.Equal(t, want, *got)
	})

	s.T().Run("IsDefault", func(t *testing.T) {
		template := test_db.GenerateEntity(func(t *test_db.Template) {
			t.IsDefault = true
			t.ProjectID = nil
			t.AuthorID = nil
		})
		templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
		require.NoError(t, err)
		defer func() { require.NoError(t, test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

		got, err := repo.GetByID(ctx, templateID)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func (s *repositorySuite) TestRepository_DeleteByID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) { p.AuthorID = userID })
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.ProjectID = &projectID
		t.AuthorID = &userID
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)

	// template version
	templateVersion := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userID
		v.Number = 1
	})
	templateVersionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", templateVersion)
	require.NoError(s.T(), err)

	// variable
	variable := test_db.GenerateEntity(func(v *test_db.Variable) {
		v.VersionID = templateVersionID
	})
	variableID, err := test_db.InsertEntityWithID[int64](s.C(), "variable", variable)
	require.NoError(s.T(), err)

	// variable constraint
	variableConstraint := test_db.GenerateEntity(func(v *test_db.Constraint) {
		v.VariableID = variableID
	})
	variableConstraintID, err := test_db.InsertEntityWithID[int64](s.C(), "variable_constraint", variableConstraint)
	require.NoError(s.T(), err)

	// handle
	err = repo.DeleteByID(ctx, templateID)
	require.NoError(s.T(), err)

	templates, err := test_db.SelectEntitiesByID[test_db.Template](s.C(), "template", []int64{templateID})
	require.NoError(s.T(), err)
	require.Empty(s.T(), templates)

	versions, err := test_db.SelectEntitiesByID[test_db.Version](s.C(), "template_version", []int64{templateVersionID})
	require.NoError(s.T(), err)
	require.Empty(s.T(), versions)

	variables, err := test_db.SelectEntitiesByID[test_db.Variable](s.C(), "variable", []int64{variableID})
	require.NoError(s.T(), err)
	require.Empty(s.T(), variables)

	variableConstraints, err := test_db.SelectEntitiesByID[test_db.Constraint](s.C(), "variable_constraint", []int64{variableConstraintID})
	require.NoError(s.T(), err)
	require.Empty(s.T(), variableConstraints)
}
