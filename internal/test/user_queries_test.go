package testdb

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_pg "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/opaulochaves/myserver/internal/entity"
	users "github.com/opaulochaves/myserver/internal/user"
	"github.com/opaulochaves/myserver/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type migration struct {
	Migrate *migrate.Migrate
}

func (this *migration) Up() (error, bool) {
	err := this.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil, true
		}
		return err, false
	}
	return nil, true
}

func (this *migration) Down() (error, bool) {
	err := this.Migrate.Down()
	if err != nil {
		return err, false
	}
	return nil, true
}

func runMigration(db *sqlx.DB, migrationsDirLocation string) (*migration, error) {
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

	return &migration{Migrate: m}, nil
}

type UserQueriesSuite struct {
	suite.Suite
	db        *sqlx.DB
	tx        *sqlx.Tx
	migration *migration
}

func (s *UserQueriesSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_URL_TEST")
	if dsn == "" {
		dsn = "postgresql://postgres:secret@localhost:5432/opchav_test"
	}

	var err error
	s.db, err = sqlx.Connect("pgx", dsn)

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	if err := s.db.Ping(); err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}

	migrationDir := "../../db/migrations"
	s.migration, err = runMigration(s.db, migrationDir)
	require.NoError(s.T(), err)

}

func (s *UserQueriesSuite) TearDownSuite() {
	// log.Println("Finishing Test. Dropping The Database")

	defer s.db.Close()
	// TODO run migration.Down only once here. Use random fake data

	// err, _ := s.migration.Down()
	// require.NoError(s.T(), err)
	// log.Println("Database Dropped Successfully")
}

func (s *UserQueriesSuite) SetupTest() {
	err, _ := s.migration.Up()
	require.NoError(s.T(), err)
	s.tx = s.db.MustBegin()
}

func (s *UserQueriesSuite) TearDownTest() {
	s.tx.Commit()

	err, _ := s.migration.Down()
	require.NoError(s.T(), err)
}

func (m *UserQueriesSuite) TestCreateUser() {
	mockUser := &getMockArrUser()[0]

	queries := users.NewUserQueries(m.db)

	res, err := queries.CreateUser(mockUser)

	require.NoError(m.T(), err)

	assert.Equal(m.T(), mockUser.ID, res.ID)
	assert.Equal(m.T(), mockUser.Email, res.Email)
	assert.Equal(m.T(), mockUser.FirstName, res.FirstName)
}

func (m *UserQueriesSuite) TestGetUser() {
	mockUser := &getMockArrUser()[0]

	queries := users.NewUserQueries(m.db)

	_, err := queries.CreateUser(mockUser)

	require.NoError(m.T(), err)

	user, err := queries.GetUser(mockUser.GetID())

	require.NoError(m.T(), err)

	assert.Equal(m.T(), mockUser.ID, user.ID)
	assert.Equal(m.T(), mockUser.Email, user.Email)
}

func (m *UserQueriesSuite) TestUpdateUser() {
	mockUser := &getMockArrUser()[0]

	queries := users.NewUserQueries(m.db)

	user, err := queries.CreateUser(mockUser)

	require.NoError(m.T(), err)

	user.FirstName = "Update First"

	userUpdated, err := queries.UpdateUser(user)

	require.NoError(m.T(), err)

	require.Equal(m.T(), userUpdated.FirstName, "Update First")
}

func (m *UserQueriesSuite) TestDeleteUser() {
	mockUser := &getMockArrUser()[0]

	queries := users.NewUserQueries(m.db)

	user, err := queries.CreateUser(mockUser)

	require.NoError(m.T(), err)

	err = queries.DeleteUser(user.ID)

	assert.Nil(m.T(), err)
}

func TestUserQueriesTestSuite(t *testing.T) {
	suite.Run(t, new(UserQueriesSuite))
}

// import (
//
//	"io"
//	"log"
//	"os"
//	"testing"
//
//	"github.com/jmoiron/sqlx"
//	"github.com/opaulochaves/myserver/internal/entity"
//	"github.com/opaulochaves/myserver/internal/user"
//	"github.com/opaulochaves/myserver/service"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
//
// )
//
//	func TestUserQueriesSuite(t *testing.T) {
//		if testing.Short() {
//			t.Skip("Skip user queries test")
//		}
//
//		dsn := os.Getenv("DATABASE_URL_TEST")
//		if dsn == "" {
//			dsn = "postgresql://postgres:secret@localhost:5432/opchav_test"
//		}
//
//		s := new(PgSuite)
//
//		suite.Run(t, s)
//	}
//
//	func (s *PgSuite) SetupTest() {
//		log.Println("Starting a test. Migrate the Database")
//		err, _ := s.Migration.Up()
//		require.NoError(s.T(), err)
//		log.Println("Database Migrated Successfully")
//		s.Tx = s.DB.MustBegin()
//	}
//
//	func (s *PgSuite) TearDownTest() {
//		log.Println("Finishing Test. Dropping The Database")
//		err, _ := s.Migration.Down()
//		require.NoError(s.T(), err)
//		log.Println("Database Dropped Successfully")
//	}
//
// // https://blevesearch.com/news/Deferred-Cleanup,-Checking-Errors,-and-Potential-Problems/
//
//	func Close(c io.Closer) {
//		err := c.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// // func getMockArrCategory() []models.Category {
// // 	return []models.Category{
// // 		models.Category{
// // 			ID:   1,
// // 			Name: "Tekno",
// // 			Slug: "tekno",
// // 		},
// // 		models.Category{
// // 			ID:   2,
// // 			Name: "Bola",
// // 			Slug: "bola",
// // 		},
// // 		models.Category{
// // 			ID:   3,
// // 			Name: "Asmara",
// // 			Slug: "asmara",
// // 		},
// // 		models.Category{
// // 			ID:   4,
// // 			Name: "Celebs",
// // 			Slug: "celebs",
// // 		},
// // 	}
// // }
func getMockArrUser() []entity.User {
	password, err := service.HashPassword("12345678")

	if err != nil {
		log.Fatal("Could not generate password")
	}

	return []entity.User{
		{
			BaseEntity: entity.BaseEntity{
				ID: 1,
			},
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  password,
		},
		{
			BaseEntity: entity.BaseEntity{
				ID: 2,
			},
			FirstName: "Mary",
			LastName:  "Doe",
			Email:     "mary@example.com",
			Password:  password,
		},
	}
}

//
// func seedUserData(t *testing.T, db *sqlx.DB) {
// 	arrUsers := getMockArrUser()
//
// 	tx := db.MustBegin()
// 	query := `INSERT INTO users (id, email, password, first_name, last_name) VALUES ($1, $2, $3, $4, $5) RETURNING *`
//
// 	for _, u := range arrUsers {
// 		_, err := tx.Exec(query, u.ID, u.Email, u.Password, u.FirstName, u.LastName)
// 		if err != nil {
// 			tx.Rollback()
// 			require.NoError(t, err)
// 		}
// 	}
//
// 	err := tx.Commit()
// 	require.NoError(t, err)
// }
//
// func (m *PgSuite) TestGetUser(t *testing.T) {
// 	mockUser := getMockArrUser()[0]
//
// 	seedUserData(m.T(), m.DB)
//
// 	queries := user.NewUserQueries(m.DB)
//
// 	res, err := queries.GetUser(mockUser.ID)
//
// 	require.NoError(m.T(), err)
//
// 	assert.Equal(m.T(), mockUser.ID, res.ID)
// 	assert.Equal(m.T(), mockUser.Email, res.Email)
// 	assert.Equal(m.T(), mockUser.FirstName, res.FirstName)
// }
//
// // func getCategoryByID(t *testing.T, DBconn *sql.DB, id int64) *models.Category {
// // 	res := &models.Category{}
// //
// // 	query := `SELECT id, name, slug, created_at, updated_at FROM category WHERE id=?`
// //
// // 	row := DBconn.QueryRow(query, id)
// // 	err := row.Scan(
// // 		&res.ID,
// // 		&res.Name,
// // 		&res.Slug,
// // 		&res.CreatedAt,
// // 		&res.UpdatedAt,
// // 	)
// // 	if err == nil {
// // 		return res
// // 	}
// // 	if err != sql.ErrNoRows {
// // 		require.NoError(t, err)
// // 	}
// // 	return nil
// // }
//
// // func getUserByID(t *testing.T, db *sql.DB)
