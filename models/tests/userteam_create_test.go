package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestUserTeamCreate tests the UserTeamCreate function
func TestUserTeamCreate(t *testing.T) {
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

	// Check if userteam associate made and expect it not to be
	var userTeam models.UserTeam
	db.Database.Where("user_id = ? AND team_id = ?", user.ID, team.ID).First(&userTeam)
	if userTeam.UserID != 0 {
		t.Fatalf("UserTeam association should not have been made")
	}

	// Call the UserTeamCreate function
	err = models.UserTeamCreate(uint64(user.ID), uint64(team.ID))
	if err != nil {
		t.Fatalf("Error creating user_team: %v", err)
	}

	// Check if userteam associate made and expect it to be
	db.Database.Where("user_id = ? AND team_id = ?", user.ID, team.ID).First(&userTeam)
	if userTeam.UserID == 0 {
		t.Fatalf("UserTeam association should have been made")
	}

	// Deferred cleanup
	defer func() {
		// Delete the user_team
		db.Database.Unscoped().Delete(&userTeam)
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
