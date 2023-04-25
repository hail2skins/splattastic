package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

// TestUserTypeShow tests the model function to show a User Type
func TestUserTypeShow(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a test user type
	testUserType := UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	// Test UserTypeShow with valid ID
	userType, err := UserTypeShow(uint64(testUserType.ID))
	assert.NoError(t, err)
	assert.NotNil(t, userType)
	assert.Equal(t, testUserType.ID, userType.ID)
	assert.Equal(t, testUserType.Name, userType.Name)

	// Test UserTypeShow with invalid ID
	_, err = UserTypeShow(0) // Assuming 0 is an invalid ID
	assert.Error(t, err)

	// Clean up
	db.Database.Unscoped().Delete(&testUserType)
}
