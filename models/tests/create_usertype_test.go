package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserType(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a user type
	userTypeName := "Test User Type"
	userType, err := models.CreateUserType(userTypeName)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the expected user type was created
	assert.Equal(t, userTypeName, userType.Name)

	// Cleanup
	db.Database.Unscoped().Delete(userType)
}
