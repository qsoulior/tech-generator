package test_db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/qsoulior/tech-generator/backend/internal/pkg/postgres"
)

type Container struct {
	db      *sqlx.DB
	builder *sq.StatementBuilderType
}

func New(db *sqlx.DB, builder *sq.StatementBuilderType) *Container {
	return &Container{db, builder}
}

func (c *Container) DB() *sqlx.DB { return c.db }

func (c *Container) Close() error { return c.db.Close() }

func NewPsql() (*Container, error) {
	db, err := postgres.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Container{db: db, builder: &builder}, nil
}

type PsqlTestSuite struct {
	suite.Suite
	container *Container
}

func (s *PsqlTestSuite) C() *Container {
	return s.container
}

func (s *PsqlTestSuite) SetupSuite() {
	container, err := NewPsql()
	require.NoError(s.T(), err)

	s.container = container
}

func (s *PsqlTestSuite) TearDownSuite() {
	require.NoError(s.T(), s.container.Close())
}
