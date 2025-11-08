package version_repository

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/version_list/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_ListByTemplateID() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	// templates
	templates := test_db.GenerateEntities(2, func(t *test_db.Template, _ int) {
		t.IsDefault = false
		t.ProjectID = nil
		t.AuthorID = &userID
	})
	templateIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template", templates)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template", templateIDs)) }()

	// template version
	versions := slices.Concat(
		test_db.GenerateEntities(3, func(v *test_db.Version, _ int) {
			v.TemplateID = templateIDs[0]
			v.AuthorID = &userID
		}),
		test_db.GenerateEntities(2, func(v *test_db.Version, _ int) {
			v.TemplateID = templateIDs[1]
			v.AuthorID = &userID
		}),
	)
	versionIDs, err := test_db.InsertEntitiesWithID[int64](s.C(), "template_version", versions)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntitiesByID(s.C(), "template_version", versionIDs)) }()

	got, err := repo.ListByTemplateID(ctx, templateIDs[0])
	require.NoError(s.T(), err)

	want := lo.Map(versions[:3], func(v test_db.Version, _ int) domain.Version {
		return domain.Version{
			ID:         v.ID,
			Number:     v.Number,
			AuthorName: user.Name,
			CreatedAt:  v.CreatedAt.Truncate(1 * time.Microsecond),
		}
	})
	slices.SortFunc(want, func(a, b domain.Version) int { return int(b.ID - a.ID) })
	require.Equal(s.T(), want, got)
}
