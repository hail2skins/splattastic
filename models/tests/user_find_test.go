package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUserFind(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	usertype := models.UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	// Create a test user
	user := models.User{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		UserType:  usertype,
	}
	err := db.Database.Create(&user).Error
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Test that the function returns the correct user
	result, err := models.UserFind(uint64(user.ID))
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)

	// Test that the function returns an error for an invalid user ID
	result, err = models.UserFind(999)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Cleanup
	db.Database.Unscoped().Delete(user)
	db.Database.Unscoped().Delete(&usertype)
}
