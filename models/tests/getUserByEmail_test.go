package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	LoadEnv()
	db.Connect()

	usertype := models.UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	email := "testemail@example.com"
	username := "testusername"
	password := "testpassword"
	firstname := "Test"
	lastname := "User"

	// Create a test user
	user, err := models.UserCreate(email, password, username, firstname, lastname, usertype.Name)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Get the user by email
	foundUser, err := models.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, email, foundUser.Email)

	// Clean up the user
	db.Database.Unscoped().Delete(user)
	db.Database.Unscoped().Delete(&usertype)
}
