package user

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opaulochaves/myserver/internal/entity"
	"github.com/opaulochaves/myserver/service"
	"github.com/pkg/errors"
)

type UserQueries interface {
	GetUsers() ([]entity.User, error)
	GetUser(id int64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int64) error
}

// userQueries struct for queries from User model.
type userQueries struct {
	db *sqlx.DB
}

func NewUserQueries(db *sqlx.DB) UserQueries {
	return &userQueries{db}
}

// BetUserByEmail implements UserQueries
func (q *userQueries) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	query := `SELECT * FROM users WHERE email = $1`

	err := q.db.Get(&user, query, email)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

// CreateUser implements UserQueries
func (q *userQueries) CreateUser(u *entity.User) (*entity.User, error) {
	// TODO: when creating a user is part of a multi step process like
	// creating user and some dependencies like categories which inserts data
	// on another table this CreateUser func must not start the transaction but
	// instead receive the transactio as param *sqlx.Tx and this tx will be used
	// here and also where categories are inserted. This runs in a service.
	// Service will be UserService or AuthService and the func that holds the tx
	// will be RegisterUser
	tx := q.db.MustBegin()
	query := `INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING *`

	var user entity.User

	hashedPassword, err := service.HashPassword(u.Password)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "hashing password error")
	}

	err = tx.QueryRowx(query, u.Email, hashedPassword, u.FirstName, u.LastName).StructScan(&user)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "insert user error")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "tx.Commit() insert user error")
	}

	return &user, nil
}

// DeleteUser implements UserQueries
func (q *userQueries) DeleteUser(id int64) error {
	// NOTE: is a tx needed for a single query?
	// TODO: remove transaction from funcs for single query
	tx := q.db.MustBegin()

	query := `DELETE FROM users WHERE id = $1`

	_, err := tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "delete user error")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "tx.Commit() delete user error")
	}

	return nil
}

// GetUser implements UserQueries
func (q *userQueries) GetUser(id int64) (*entity.User, error) {
	var user entity.User

	query := `SELECT * FROM users WHERE id = $1`

	err := q.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUsers implements UserQueries
func (q *userQueries) GetUsers() ([]entity.User, error) {
	// users := []entity.User{}
	var users []entity.User

	query := `SELECT * FROM users`

	err := q.db.SelectContext(context.Background(), &users, query)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser implements UserQueries
func (q *userQueries) UpdateUser(u *entity.User) (*entity.User, error) {
	// NOTE: is a tx needed for a single query?
	tx := q.db.MustBegin()

	query := `UPDATE users SET first_name = $2, last_name = $3, updated_at = $4 WHERE id = $1 RETURNING *`

	var user entity.User

	err := tx.QueryRowx(query, u.ID, u.FirstName, u.LastName, time.Now()).StructScan(&user)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "update user error")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "tx.Commit() update user error")
	}

	return &user, nil
}
