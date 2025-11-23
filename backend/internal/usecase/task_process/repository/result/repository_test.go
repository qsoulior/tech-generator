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

	data := []byte{1, 2, 3}
	resultID, err := repo.Insert(ctx, data)
	require.NoError(s.T(), err)
	defer func() { require.NoError(s.T(), test_db.DeleteEntityByID(s.C(), "result", resultID)) }()

	gotResults, err := test_db.SelectEntitiesByID[test_db.Result](s.C(), "result", []int64{resultID})
	require.NoError(s.T(), err)
	require.Len(s.T(), gotResults, 1)
	require.Equal(s.T(), data, gotResults[0].Data)
}
