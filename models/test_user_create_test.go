package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUserCreate(t *testing.T) {
	LoadEnv()
	db.Connect()

	usertype := UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	// Test data
	email := "test@example.com"
	password := "test_password"
	username := "test_username"
	firstname := "Test"
	lastname := "User"
	usertypeName := usertype.Name

	// Call UserCreate function
	user, err := UserCreate(email, password, username, firstname, lastname, usertypeName)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the created user has the expected values
	assert.NotNil(t, user)
	assert.NotZero(t, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, username, user.UserName)
	assert.Equal(t, firstname, user.FirstName)
	assert.Equal(t, lastname, user.LastName)
	assert.Equal(t, usertype.ID, user.UserTypeID)

	// Check if the password is hashed
	assert.NotEqual(t, password, user.Password)
	assert.True(t, helpers.CheckPasswordHash(password, user.Password))

	// Cleanup
	db.Database.Unscoped().Delete(user)
	db.Database.Unscoped().Delete(&usertype)
}
