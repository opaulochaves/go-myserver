package test

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TSuite struct {
	suite.Suite
	DB             *sqlx.DB
	TX             *sqlx.Tx
	Migration      *MigrationTest
	TruncateTables string
}

func (t *TSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_URL_TEST")
	if dsn == "" {
		dsn = "postgresql://postgres:secret@localhost:5432/opchav_test"
	}

	var err error
	t.DB, err = sqlx.Connect("pgx", dsn)

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	if err := t.DB.Ping(); err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}

	t.TruncateTables = "notes, users"
	t.Migration, err = runMigration(t.DB, "../../db/migrations")

	require.NoError(t.T(), err)
}

func (t *TSuite) SetupTest() {
	err, _ := t.Migration.Up()
	require.NoError(t.T(), err)
	t.TX = t.DB.MustBegin()
}

func (t *TSuite) TearDownTest() {
	// TODO: handle rollback
	err := t.TX.Commit()
	require.NoError(t.T(), err)

	query := fmt.Sprintf("TRUNCATE %s RESTART IDENTITY;", t.TruncateTables)
	_, err = t.DB.Exec(query)

	require.NoError(t.T(), err)
}

func (t *TSuite) TearDownSuite() {
	defer t.DB.Close()

	err, _ := t.Migration.Down()
	require.NoError(t.T(), err)
}
