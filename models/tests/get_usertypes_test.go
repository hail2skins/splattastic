package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserTypes(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create user types
	userType1 := models.UserType{Name: "Test User Type"}
	userType2 := models.UserType{Name: "Another User Type"}
	db.Database.Create(&userType1)
	db.Database.Create(&userType2)

	// Call GetUserTypes function
	userTypes, err := models.GetUserTypes()

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the expected user types were returned
	found1 := false
	found2 := false
	for _, userType := range userTypes {
		if userType.Name == userType1.Name {
			found1 = true
		}
		if userType.Name == userType2.Name {
			found2 = true
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)

	// Cleanup
	db.Database.Unscoped().Delete(&userType1)
	db.Database.Unscoped().Delete(&userType2)
}
