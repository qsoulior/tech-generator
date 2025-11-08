package version_repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
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

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// template
	template := test_db.GenerateEntity(func(t *test_db.Template) {
		t.IsDefault = false
		t.ProjectID = nil
		t.AuthorID = &userID
	})
	templateID, err := test_db.InsertEntityWithID[int64](s.C(), "template", template)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template", templateID)) }()

	// template version
	want := test_db.GenerateEntity(func(v *test_db.Version) {
		v.TemplateID = templateID
		v.AuthorID = &userID
		v.Number = 1
	})

	templateVersion := domain.Version{
		TemplateID: templateID,
		AuthorID:   userID,
		Data:       want.Data,
	}

	templateVersionID, err := repo.Create(ctx, templateVersion)
	require.NoError(s.T(), err)
	defer func() {
		require.NoError(s.T(), test_db.DeleteEntityByColumn(s.C(), "template_version", "template_id", templateID))
	}()

	templates, err := test_db.SelectEntitiesByID[test_db.Version](s.C(), "template_version", []int64{templateVersionID})
	require.NoError(s.T(), err)
	require.Len(s.T(), templates, 1)

	got := templates[0]
	want.ID = templateVersionID
	want.CreatedAt = got.CreatedAt
	require.Equal(s.T(), want, got)
}
