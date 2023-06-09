package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUserTypeDelete(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Create a test user type
	testUserType := models.UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	// Soft delete the test user type
	err := models.UserTypeDelete(uint64(testUserType.ID))
	assert.NoError(t, err)

	// Check if the user type is soft deleted in the database
	var softDeletedUserType models.UserType
	result := db.Database.Unscoped().First(&softDeletedUserType, testUserType.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, softDeletedUserType.DeletedAt)

	// Cleanup
	db.Database.Unscoped().Delete(&testUserType)
}
