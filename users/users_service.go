package users

import (
	"log"
	"strings"

	"github.com/opaulochaves/myserver/apperrors"
	"github.com/opaulochaves/myserver/service"
)

// UserService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository UserRepository
}

// USConfig will hold repositories that will eventually be injected into
// this service layer
type USConfig struct {
	UserRepository UserRepository
}

// UserService defines methods related to account operations the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) error
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

// CreateUser implements UserService
func (s *userService) CreateUser(user *User) (*User, error) {
	hashedPassword, err := service.HashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}

	user.Password = hashedPassword

	return s.UserRepository.Create(user)
}

// Get implements UserService
func (s *userService) Get(id int) (*User, error) {
	return s.UserRepository.FindByID(id)
}

// GetByEmail implements UserService
func (s *userService) GetByEmail(email string) (*User, error) {

	// Sanitize email
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)

	return s.UserRepository.FindByEmail(email)
}

// UpdateUser implements UserService
func (s *userService) UpdateUser(user *User) error {
	return s.UserRepository.Update(user)
}
