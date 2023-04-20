package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

func TestCheckEmailUsernameAvailable(t *testing.T) {
	LoadEnv()
	db.Connect()
	// Create a user for testing
	user, err := UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"Test",
		"User",
		Athlete,
	)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Check if the email and username are available (should be false)
	available := CheckEmailUsernameAvailable("test@example.com", "testuser")
	assert.False(t, available)

	// Cleanup
	db.Database.Unscoped().Delete(user)
}
