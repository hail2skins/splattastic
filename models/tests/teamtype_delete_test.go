package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamTypeDelete tests the TeamTypeDelete function
func TestTeamTypeDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new team_type
	team_type, err := models.TeamTypeCreate("TestTeamType")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the team_type
	err = models.TeamTypeDelete(uint64(team_type.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Verify the team_type was deleted
	_, err = models.TeamTypeShow(uint64(team_type.ID))
	if err == nil {
		t.Errorf("TeamType with ID %d not deleted", team_type.ID)
	}

	// Delete the team_type
	db.Database.Unscoped().Delete(&team_type)
}
