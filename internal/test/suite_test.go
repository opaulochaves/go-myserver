package testdb
//
// import (
// 	"log"
// 	"os"
//
// 	_ "github.com/jackc/pgx/v5/stdlib"
// 	"github.com/jmoiron/sqlx"
//
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// )
//
// const pg = "postgres"
//
// type PgSuite struct {
// 	suite.Suite
// 	DB        *sqlx.DB
// 	Tx        *sqlx.Tx
// 	Migration *migration
// }
//
// func (s *PgSuite) SetupSuite() {
// 	var err error
//
// 	dsn := os.Getenv("DATABASE_URL_TEST")
// 	if dsn == "" {
// 		dsn = "postgresql://postgres:secret@localhost:5432/opchav_test"
// 	}
//
// 	db, err := sqlx.Connect("pgx", dsn)
//
// 	if err != nil {
// 		log.Fatalf("error, not connected to database, %v", err)
// 	}
//
// 	for {
// 		err = db.Ping()
// 		if err == nil {
// 			break
// 		}
// 	}
//
// 	migrationDir := "db/migrations"
// 	s.Migration, err = runMigration(s.DB.DB, migrationDir)
//
// 	require.NoError(s.T(), err)
// }
//
// // TearDownSuite teardown at the end of test
// func (s *PgSuite) TearDownSuite() {
// 	log.Println("Finishing Test. Dropping The Database")
//
// 	defer s.DB.Close()
//
// 	err, _ := s.Migration.Down()
// 	require.NoError(s.T(), err)
// 	log.Println("Database Dropped Successfully")
// }
