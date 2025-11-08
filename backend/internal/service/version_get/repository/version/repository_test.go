package version_repository

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
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

func (s *repositorySuite) TestRepository_GetByID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	s.T().Run("Exists", func(t *testing.T) {
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
		templateVersion := test_db.GenerateEntity(func(v *test_db.TemplateVersion) {
			v.TemplateID = templateID
			v.AuthorID = &userID
		})
		templateVersionID, err := test_db.InsertEntityWithID[int64](s.C(), "template_version", templateVersion)
		require.NoError(s.T(), err)
		defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "template_version", templateVersionID)) }()

		got, err := repo.GetByID(ctx, templateVersionID)
		require.NoError(t, err)

		want := domain.Version{
			ID:         templateVersionID,
			TemplateID: templateID,
			Number:     templateVersion.Number,
			CreatedAt:  templateVersion.CreatedAt.Truncate(1 * time.Microsecond),
			Data:       templateVersion.Data,
		}
		require.Equal(t, want, *got)
	})

	s.T().Run("NotExists", func(t *testing.T) {
		got, err := repo.GetByID(ctx, gofakeit.Int64())
		require.NoError(t, err)
		require.Nil(t, got)
	})
}
