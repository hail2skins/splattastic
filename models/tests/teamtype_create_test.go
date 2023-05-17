package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamTypeCreate tests the TeamTypeCreate function
func TestTeamTypeCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a team type
	teamType, err := models.TeamTypeCreate("TestTeamType")
	if err != nil {
		t.Errorf("Error creating team type: %v", err)
	}

	// Test that the team type was created
	if teamType.Name != "TestTeamType" {
		t.Errorf("TeamTypeCreate returned an incorrect name: got %v want %v", teamType.Name, "TestTeamType")
	}

	// Test duplicate team type name
	dupTeamType, err := models.TeamTypeCreate("TestTeamType")
	if err == nil {
		t.Errorf("Expected error when creating team type with duplicate name")
	}

	// Test blank team type name
	blankTeamType, err := models.TeamTypeCreate("")
	if err == nil {
		t.Errorf("Expected error when creating team type with blank name")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&teamType)
	db.Database.Unscoped().Delete(&dupTeamType)
	db.Database.Unscoped().Delete(&blankTeamType)
}
