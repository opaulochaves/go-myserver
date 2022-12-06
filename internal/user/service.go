package user

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/opaulochaves/myserver/internal/entity"
	"github.com/opaulochaves/myserver/internal/util"
	"github.com/pkg/errors"
)

type Service interface {
	Get(id int64) (User, error)
	Query(offset int, limit int) ([]User, error)
	Count() (int, error)
	Create(input CreateUserRequest) (User, error)
	Update(id int64, input UpdateUserRequest) (User, error)
	Delete(id int64) (User, error)
}

// User represents the data about an user.
type User struct {
	*entity.User
}

// CreateUserRequest represents an user creation request.
type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Bind implements render.Binder
func (*CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

type UserResponse struct {
	User
}

// Render implements render.Renderer
func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Validate validates the CreateUserRequest fields.
func (c CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.LastName, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 100)),
	)
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// Validate validates the CreateUserRequest fields.
func (c UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Length(2, 255)),
		validation.Field(&c.LastName, validation.Length(2, 255)),
		validation.Field(&c.Password, validation.Length(8, 100)),
	)
}

type service struct {
	repo UserQueries
	// logger log.Logger
}

func NewService(repo UserQueries) Service {
	return service{repo}
}

// Count implements Service
func (s service) Count() (int, error) {
	return s.repo.Count()
}

// Create implements Service
func (s service) Create(input CreateUserRequest) (User, error) {
	if err := input.Validate(); err != nil {
		return User{}, err
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return User{}, errors.Wrap(err, "hashing password error")
	}

	user, err := s.repo.CreateUser(&entity.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		return User{}, err
	}

	return User{user}, nil
}

// Get implements Service
func (s service) Get(id int64) (User, error) {
	user, err := s.repo.GetUser(id)

	return User{user}, err
}

const defaultLimit = 10

// Query implements Service
func (s service) Query(offset int, limit int) ([]User, error) {
	users, err := s.repo.GetUsers(offset, limit)
	if err != nil {
		return nil, err
	}

	var result []User

	for _, u := range users {
		// NOTE: I still don't understand that fully
		// https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		u := u
		result = append(result, User{User: &u})
	}

	return result, nil
}

// Update implements Service
func (s service) Update(id int64, input UpdateUserRequest) (User, error) {
	if err := input.Validate(); err != nil {
		return User{}, err
	}

	var err error
	var hashedPassword string

	if input.Password != "" {
		hashedPassword, err = util.HashPassword(input.Password)
	}

	if err != nil {
		return User{}, errors.Wrap(err, "hashing password error")
	}

	user, err := s.Get(id)
	if err != nil {
		return user, err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName

	if hashedPassword != "" {
		user.Password = hashedPassword
	}

	if _, err := s.repo.UpdateUser(user.User); err != nil {
		return user, err
	}

	return user, err
}

// Delete implements Service
func (s service) Delete(id int64) (User, error) {
	user, err := s.Get(id)
	if err != nil {
		return User{}, err
	}

	if err = s.repo.DeleteUser(id); err != nil {
		return User{}, nil
	}

	return user, nil
}
