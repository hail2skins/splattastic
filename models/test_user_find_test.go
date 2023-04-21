package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

func TestUserFind(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Create a test user
	// Create a test user
	user, err := UserCreate("test@example.com", "testpassword", "John", "Doe", "testuser", UserType("Athlete"))
	if err != nil {
		t.Errorf("Failed to create test user: %v", err)
	}

	// Test that the function returns the correct user
	result, err := UserFind(uint64(user.ID))
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)

	// Test that the function returns an error for an invalid user ID
	result, err = UserFind(999)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Cleanup
	db.Database.Unscoped().Delete(user)
}
