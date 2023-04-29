package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestUserUpdate tests the user update functionality
func TestUserUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two user types
	ut1, _ := models.CreateUserType("Test User Type 1")
	ut2, _ := models.CreateUserType("Test User Type 2")

	// Create a user
	email := "test@example.com"
	password := "testpassword"
	firstName := "John"
	lastName := "Doe"
	userName := "testuser"
	usertypeName := ut1.Name

	user, err := models.UserCreate(email, password, firstName, lastName, userName, usertypeName)

	// Update the user
	newEmail := "test1@example.com"
	newFirstName := "Jane"
	newLastName := "Dime"
	newUserName := "testuser1"
	newUserTypeID := ut2.ID

	err = user.Update(newEmail, newFirstName, newLastName, newUserName, uint64(newUserTypeID))
	if err != nil {
		t.Errorf("error updating user: %v", err)
	}

	// Fetch the updated user
	updatedUser, err := models.UserShow(uint64(user.ID))
	if err != nil {
		t.Errorf("error showing user: %v", err)
	}

	// Check if the user has been updated
	if updatedUser.Email != newEmail ||
		updatedUser.FirstName != newFirstName ||
		updatedUser.LastName != newLastName ||
		updatedUser.UserName != newUserName ||
		updatedUser.UserType.ID != newUserTypeID {
		t.Errorf("user has not been updated")
	}

	// Clean up
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&ut1)
	db.Database.Unscoped().Delete(&ut2)

}
