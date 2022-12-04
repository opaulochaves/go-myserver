package user

import (
	"testing"

	"github.com/opaulochaves/myserver/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type queriesSuiteTest struct {
	test.TSuite
}

func TestQueriesSuiteTest(t *testing.T) {
	suite.Run(t, new(queriesSuiteTest))
}

func (t *queriesSuiteTest) TestCreateUser() {
	mockUser := &test.GenerateUsers(1)[0]

	queries := NewUserQueries(t.DB, t.TX)

	res, err := queries.CreateUser(mockUser)

	require.NoError(t.T(), err)

	assert.Equal(t.T(), mockUser.Email, res.Email)
	assert.Equal(t.T(), mockUser.FirstName, res.FirstName)
}

func (t *queriesSuiteTest) TestGetUser() {
	mockUser := &test.GenerateUsers(1)[0]

	queries := NewUserQueries(t.DB, t.TX)

	userSaved, err := queries.CreateUser(mockUser)

	require.NoError(t.T(), err)

	user, err := queries.GetUser(userSaved.ID)

	require.NoError(t.T(), err)

	assert.Equal(t.T(), userSaved.ID, user.ID)
	assert.Equal(t.T(), userSaved.Email, user.Email)
}

func (t *queriesSuiteTest) TestGetUserByEmail() {
	mockUser := &test.GenerateUsers(1)[0]

	queries := NewUserQueries(t.DB, t.TX)

	_, err := queries.CreateUser(mockUser)

	require.NoError(t.T(), err)

	user, err := queries.GetUserByEmail(mockUser.Email)

	require.NoError(t.T(), err)

	assert.Equal(t.T(), mockUser.Email, user.Email)
}

func (t *queriesSuiteTest) TestUpdateUser() {
	mockUser := &test.GenerateUsers(1)[0]

	queries := NewUserQueries(t.DB, t.TX)

	user, err := queries.CreateUser(mockUser)

	require.NoError(t.T(), err)

	user.FirstName = "Update First"

	userUpdated, err := queries.UpdateUser(user)

	require.NoError(t.T(), err)

	require.Equal(t.T(), userUpdated.FirstName, "Update First")
}

func (t *queriesSuiteTest) TestDeleteUser() {
	mockUser := &test.GenerateUsers(1)[0]

	queries := NewUserQueries(t.DB, t.TX)

	user, err := queries.CreateUser(mockUser)

	require.NoError(t.T(), err)

	err = queries.DeleteUser(user.ID)

	assert.Nil(t.T(), err)
}
