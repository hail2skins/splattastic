package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestUserShow tests the UserShow function
func TestUserShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a test user type
	usertype := models.UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	// Create a test user
	user := models.User{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		UserType:  usertype,
	}
	err := db.Database.Create(&user).Error
	if err != nil {
		t.Fatal(err)
	}

	// Get the user
	userToShow, err := models.UserShow(uint64(user.ID))
	if err != nil {
		t.Errorf("error showing user: %v", err)
	}
	// Nil check
	if userToShow == nil {
		t.Errorf("handler returned nil user")
	}
	// Check the user ID
	if userToShow.ID != user.ID {
		t.Errorf("handler returned unexpected user ID: got %v want %v", user.ID, userToShow.ID)
	}
	// Check the user email
	if userToShow.Email != "test@example.com" {
		t.Errorf("handler returned unexpected user email: got %v want %v", user.Email, "test@example.com")
	}
	// Check the user username
	if userToShow.UserName != "testuser" {
		t.Errorf("handler returned unexpected user username: got %v want %v", user.UserName, "testuser")
	}
	// Check the user first name
	if userToShow.FirstName != "John" {
		t.Errorf("handler returned unexpected user first name: got %v want %v", user.FirstName, "John")
	}
	// Check the user last name
	if userToShow.LastName != "Doe" {
		t.Errorf("handler returned unexpected user last name: got %v want %v", user.LastName, "Doe")
	}
	// Check the user type name
	if userToShow.UserType.Name != "Test User Type" {
		t.Errorf("handler returned unexpected user type name: got %v want %v", user.UserType.Name, "Test User Type")
	}

	// Delete the test user
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&usertype)

}
