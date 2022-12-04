package test

import (
	"fmt"
	"log"

	"github.com/opaulochaves/myserver/internal/entity"
	"github.com/opaulochaves/myserver/internal/util"
)

func GenerateUsers(count int) []entity.User {
	var users []entity.User

	password, err := util.HashPassword("12345678")

	if err != nil {
		log.Fatalf("hasing password error: %v", err)
	}

	for i := 0; i < count; i++ {
		id := i + 1
		user := entity.User{
			Email:     fmt.Sprintf("user0%d@example.com", id),
			FirstName: fmt.Sprintf("User 0%d", id),
			LastName:  "Example",
			Password:  password,
		}

		users = append(users, user)
	}

	return users
}
