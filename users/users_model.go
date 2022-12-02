package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents the user of the website.
type User struct {
	ID        int64              `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Password  string             `json:"-"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
} //@name User

// userRepository is data/repository implementation
// of service layer UserRepository
type userRepository struct {
	db *pgxpool.Pool
}

// UserRepository defines methods related to account db operations the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	Create(user *User) (*User, error)
	FindByID(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create implements UserRepository
func (r *userRepository) Create(user *User) (*User, error) {
	row := r.db.QueryRow(context.Background(), sqlInsertUser, user.Email, user.Password, user.FirstName, user.LastName)
	var u = &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}

// FindByEmail implements UserRepository
func (r *userRepository) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow(context.Background(), sqlSelectByEmail, email)

	var u = &User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Password, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}

// FindByID implements UserRepository
func (r *userRepository) FindByID(id int) (*User, error) {
	row := r.db.QueryRow(context.Background(), sqlSelectByID, id)

	var u = &User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Password, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}

// Update implements UserRepository
func (r *userRepository) Update(user *User) error {
	rows, err := r.db.Query(context.Background(), sqlUpdateUser, user.ID)

	defer rows.Close()

	return err
}

const sqlInsertUser = `INSERT INTO users (email, password, first_name, last_name)
	VALUES ($1, $2, $3, $4) RETURNING id, email, password, first_name, last_name, created_at, updated_at`

const sqlUpdateUser = `UPDATE users SET first_name = $1, last_name = $2, password = $3, update_at = $4 WHERE id = $5;`

const sqlSelectByID = `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users where id = $1`

const sqlSelectByEmail = `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users where email = $1`
