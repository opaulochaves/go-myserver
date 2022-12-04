package testdb

// FROM: https://github.com/bxcodec/integration-testing/blob/master/mysql/migration_test.go

// import (
// 	"database/sql"
// 	"strings"
//
// 	"github.com/golang-migrate/migrate/v4"
// 	_pg "github.com/golang-migrate/migrate/v4/database/pgx"
// 	_ "github.com/jackc/pgx/v5/stdlib"
// )
//
// type migration struct {
// 	Migrate *migrate.Migrate
// }
//
// func (this *migration) Up() (error, bool) {
// 	err := this.Migrate.Up()
// 	if err != nil {
// 		if err == migrate.ErrNoChange {
// 			return nil, true
// 		}
// 		return err, false
// 	}
// 	return nil, true
// }
//
// func (this *migration) Down() (error, bool) {
// 	err := this.Migrate.Down()
// 	if err != nil {
// 		return err, false
// 	}
// 	return nil, true
// }
//
// func runMigration(db *sql.DB, migrationsDirLocation string) (*migration, error) {
// 	dataPath := []string{}
// 	dataPath = append(dataPath, "file://")
// 	dataPath = append(dataPath, migrationsDirLocation)
//
// 	pathToMigrate := strings.Join(dataPath, "")
//
// 	driver, err := _pg.WithInstance(db, &_pg.Config{})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "postgres", driver)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &migration{Migrate: m}, nil
// }
