package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserType(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Create a test user type
	testUserType := models.UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	// Update the name of the test user type
	newName := "UpdatedTestType"
	err := testUserType.Update(newName)
	assert.NoError(t, err)

	// Retrieve the updated user type from the database
	var updatedUserType models.UserType
	db.Database.First(&updatedUserType, testUserType.ID)

	// Check if the name was updated correctly
	assert.Equal(t, newName, updatedUserType.Name)

	// Cleanup
	db.Database.Unscoped().Delete(&testUserType)
}
