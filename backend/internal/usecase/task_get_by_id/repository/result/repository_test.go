package result_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	test_db "github.com/qsoulior/tech-generator/backend/internal/pkg/test/db"
)

type repositorySuite struct {
	test_db.PsqlTestSuite
}

func Test_repositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) TestRepository_Insert() {
	ctx := context.Background()
	repo := New(s.C().DB())

	result := test_db.GenerateEntity[test_db.Result]()
	resultID, err := test_db.InsertEntityWithID[int64](s.C(), "result", result)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "result", resultID)) }()

	got, err := repo.GetDataByID(ctx, resultID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), result.Data, got)
}
