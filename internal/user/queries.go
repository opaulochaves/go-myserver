package user

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opaulochaves/myserver/internal/entity"
	"github.com/opaulochaves/myserver/internal/util"
	"github.com/pkg/errors"
)

type UserQueries interface {
	GetUsers(offset, limit int) ([]entity.User, error)
	GetUser(id int64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int64) error
	Count() (int, error)
}

// userQueries struct for queries from User model.
type userQueries struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewUserQueries(db *sqlx.DB, tx *sqlx.Tx) UserQueries {
	return &userQueries{db, tx}
}

// BetUserByEmail implements UserQueries
func (q *userQueries) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	query := `SELECT * FROM users WHERE email = $1`

	err := q.tx.Get(&user, query, email)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

// CreateUser implements UserQueries
func (q *userQueries) CreateUser(u *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING *`

	var user entity.User

	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hashing password error")
	}

	err = q.tx.QueryRowx(query, u.Email, hashedPassword, u.FirstName, u.LastName).StructScan(&user)
	if err != nil {
		return nil, errors.Wrap(err, "insert user error")
	}

	return &user, nil
}

// DeleteUser implements UserQueries
func (q *userQueries) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := q.tx.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "delete user error")
	}

	return nil
}

// GetUser implements UserQueries
func (q *userQueries) GetUser(id int64) (*entity.User, error) {
	var user entity.User

	query := `SELECT * FROM users WHERE id = $1`

	// TODO: improve what to use db? or tx? when? pass as param?
	// right now it means I must call this function GetUser within a transaction.
	// I don't need to always get the data within a transaction
	err := q.db.Get(&user, query, id)

	fmt.Println(err)
	// TODO throw not found err if not user for the given id

	return &user, err
}

// GetUsers implements UserQueries
func (q *userQueries) GetUsers(offset, limit int) ([]entity.User, error) {
	users := []entity.User{}

	query := `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`

	err := q.db.Select(&users, query, limit, offset)

	return users, err
}

// UpdateUser implements UserQueries
func (q *userQueries) UpdateUser(u *entity.User) (*entity.User, error) {
	query := `UPDATE users SET first_name = $2, last_name = $3, updated_at = $4 WHERE id = $1 RETURNING *`

	var user entity.User

	err := q.tx.QueryRowx(query, u.ID, u.FirstName, u.LastName, time.Now()).StructScan(&user)
	if err != nil {
		return nil, errors.Wrap(err, "update user error")
	}

	return &user, nil
}

// TODO: use a criteria if present to count
// Count returns the number of rows on the users table
func (q *userQueries) Count() (int, error) {
	var count int
	query := `SELECT COUNT(id) FROM users`
	err := q.db.QueryRow(query).Scan(&count)
	return count, err
}
