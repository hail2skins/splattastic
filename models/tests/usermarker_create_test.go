package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestUserMarkerCreate tests the creating a User with criteria to associate a marker
func TestUserMarkerCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a user type
	ut1, _ := models.CreateUserType("Test User Type 1")

	// Create a marker for the test
	marker, _ := models.MarkerCreate("Test Marker", "Test Marker Description")

	// Create the user
	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)

	// Fetch the user from the database to get the associated markers
	var fetchedUser models.User
	if err := db.Database.Preload("Markers").First(&fetchedUser, user.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve user with markers: %v", err)
	}

	// Check that the user has the Test Marker
	hasTestMarker := false
	for _, marker := range fetchedUser.Markers {
		if marker.Name == "Test Marker" {
			hasTestMarker = true
			break
		}
	}

	if !hasTestMarker {
		t.Fatalf("User did not have the Test Marker")
	}

	// Cleanup

	// Delete the UserMarker
	var userMarker models.UserMarker
	db.Database.Unscoped().Where("user_id = ?", user.ID).Delete(&userMarker)

	// Cleanup
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&ut1)
	db.Database.Unscoped().Delete(&marker)

}
