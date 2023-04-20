package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

func TestCheckEmailUsernameAvailable(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a test user
	user := User{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		UserType:  "Athlete",
		Admin:     false,
	}
	err := db.Database.Create(&user).Error
	assert.NoError(t, err)

	// Test with an unavailable email and username
	available, err := CheckEmailUsernameAvailable("test@example.com", "differentusername")
	assert.NoError(t, err)
	assert.False(t, available, "Email should not be available")

	// Test with an unavailable email and username
	available, err = CheckEmailUsernameAvailable("differentemail@example.com", "testuser")
	assert.NoError(t, err)
	assert.False(t, available, "Username should not be available")

	// Test with an unavailable email and username
	available, err = CheckEmailUsernameAvailable("test@example.com", "testuser")
	assert.NoError(t, err)
	assert.False(t, available, "Email and username should not be available")

	// Test with available email and username
	available, err = CheckEmailUsernameAvailable("differentemail@example.com", "differentusername")
	assert.NoError(t, err)
	assert.True(t, available, "Email and username should be available")
}
