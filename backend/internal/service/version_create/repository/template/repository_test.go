package template_repository

import (
	"context"
	"testing"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/brianvoe/gofakeit/v7"
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

func (s *repositorySuite) TestRepository_UpdateByID() {
	ctx := context.Background()
	repo := New(s.C().DB(), trmsqlx.DefaultCtxGetter)

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// project
	project := test_db.GenerateEntity(func(p *test_db.Project) {
		p.AuthorID = userID
	})
	projectID, err := test_db.InsertEntityWithID[int64](s.C(), "project", project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "project", projectID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.ProjectID = &projectID
		t.AuthorID = &userID
		t.CreatedAt = gofakeit.Date().Truncate(1 * time.Second)
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template updating
	templateToUpdate := domain.TemplateToUpdate{
		ID:            templateID,
		LastVersionID: gofakeit.Int64(),
	}
	err = repo.UpdateByID(ctx, templateToUpdate)
	require.NoError(s.T(), err)

	templates, err := test_db.SelectEntitiesByID[test_db.Template](s.C(), "template", []int64{templateID})
	require.NoError(s.T(), err)
	require.Len(s.T(), templates, 1)

	got := templates[0]
	require.Equal(s.T(), templateToUpdate.LastVersionID, *got.LastVersionID)

	now := time.Now().UTC().Truncate(1 * time.Second)
	require.GreaterOrEqual(s.T(), *got.UpdatedAt, now)
}
