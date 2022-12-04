package test

import (
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_pg "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type MigrationTest struct {
	Migrate *migrate.Migrate
}

func (m *MigrationTest) Up() (error, bool) {
	err := m.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil, true
		}
		return err, false
	}
	return nil, true
}

func (m *MigrationTest) Down() (error, bool) {
	err := m.Migrate.Down()
	if err != nil {
		return err, false
	}
	return nil, true
}

func runMigration(db *sqlx.DB, migrationsDirLocation string) (*MigrationTest, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsDirLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := _pg.WithInstance(db.DB, &_pg.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return &MigrationTest{Migrate: m}, nil
}
