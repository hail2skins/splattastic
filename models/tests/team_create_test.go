package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamCreate tests the TeamCreate function
func TestTeamCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a user type
	ut1, _ := models.CreateUserType("Test User Type")

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)

	// Create a team type
	tt1, _ := models.TeamTypeCreate("Test Team Type")

	// Create a state
	state, _ := models.StateCreate("Test State", "TS")

	// Create a team with required fields
	team, err := models.TeamCreate("Test Team", "", "", "", "12345", "", uint64(tt1.ID), uint64(state.ID))
	if err != nil {
		t.Fatalf("Error creating team: %v", err)
	}

	// Create a team without state and confirm it fails
	_, err = models.TeamCreate("Test Team", "", "", "", "12345", "", uint64(tt1.ID), 0)
	if err == nil {
		t.Fatalf("Error creating team: %v", err)
	}

	// Create a team without zip and confirm it fails
	_, err = models.TeamCreate("Test Team", "", "", "", "", "", uint64(tt1.ID), uint64(state.ID))
	if err == nil {
		t.Fatalf("Error creating team: %v", err)
	}

	// Deferred cleanup
	defer func() {
		// Delete the team
		db.Database.Unscoped().Delete(&team)
		// Delete the UserMarker
		var userMarker models.UserMarker
		db.Database.Unscoped().Where("user_id = ?", user.ID).Delete(&userMarker)
		// Delete the user
		db.Database.Unscoped().Delete(&user)
		// Delete the team type
		db.Database.Unscoped().Delete(&tt1)
		// Delete the state
		db.Database.Unscoped().Delete(&state)
		// Delete the user type
		db.Database.Unscoped().Delete(&ut1)
	}()
}
