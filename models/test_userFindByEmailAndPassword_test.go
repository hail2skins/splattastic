package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

func TestUserFindByEmailAndPassword(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a test user
	// Create a test user
	testUser, err := UserCreate("test@example.com", "testpassword", "John", "Doe", "testuser", UserType("Athlete"))
	if err != nil {
		t.Errorf("Failed to create test user: %v", err)
	}

	// Test valid email and password
	user, err := UserFindByEmailAndPassword("test@example.com", "testpassword")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.ID, testUser.ID)

	// Test invalid email
	user, err = UserFindByEmailAndPassword("invalid@example.com", "testpassword")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test invalid password
	user, err = UserFindByEmailAndPassword("test@example.com", "invalidpassword")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Cleanup
	db.Database.Unscoped().Delete(testUser)
}
