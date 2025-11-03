package project_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/project_create/domain"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_Create() {
	ctx := context.Background()
	repo := New(s.C().DB())

	// user
	user := test_db.GenerateEntity[test_db.User]()
	userID, err := test_db.InsertEntityWithID[int64](s.C(), "usr", user)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "usr", userID)) }()

	want := test_db.Project{
		Name:     gofakeit.UUID(),
		AuthorID: userID,
	}

	project := domain.Project{
		Name:     want.Name,
		AuthorID: userID,
	}
	err = repo.Create(ctx, project)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByColumn(s.C(), "project", "name", project.Name)) }()

	projects, err := test_db.SelectEntitiesByColumn[test_db.Project](s.C(), "project", "name", []string{project.Name})
	require.NoError(s.T(), err)
	require.Len(s.T(), projects, 1)

	got := projects[0]
	want.ID = got.ID
	require.Equal(s.T(), want, got)
}
