package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamTypeShow tests the TeamTypeShow function
func TestTeamTypeShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a team type
	teamType, err := models.TeamTypeCreate("TestTeamType")
	if err != nil {
		t.Fatal(err)
	}

	// Get the team type
	teamType, err = models.TeamTypeShow(uint64(teamType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the team type name
	if teamType.Name != "TestTeamType" {
		t.Errorf("handler returned unexpected team type name: got %v want %v", teamType.Name, "TestTeamType")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&teamType)
}
