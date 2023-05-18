package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamTypeUpdate is a controller for testing the update of a team_type
func TestTeamTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new team_type
	tt1, _ := models.TeamTypeCreate("Test Team Type")

	// Update the team_type name
	newName := "UpdatedTestTeamType"
	err := tt1.Update(newName)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated team_type
	updatedTeamType, err := models.TeamTypeShow(uint64(tt1.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedTeamType.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedTeamType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedTeamType)
}
