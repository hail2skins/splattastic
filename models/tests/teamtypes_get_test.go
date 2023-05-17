package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/hail2skins/splattastic/models/helpers"
)

// TestTeamTypesGet is a function which will test the TeamTypesGet function
func TestTeamTypesGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two team types
	teamType1, _ := models.TeamTypeCreate("TestTeamType1")
	teamType2, _ := models.TeamTypeCreate("TestTeamType2")

	// Get all team types
	teamTypes, err := models.TeamTypesGet()
	if err != nil {
		t.Errorf("Error getting team types: %v", err)
	}

	// Convert the team types slice to a slice of interface{}
	teamTypesInterface := make([]interface{}, len(teamTypes))
	for i, tt := range teamTypes {
		teamTypesInterface[i] = tt
	}

	// Check that the team types are the ones we created using containsModel function
	if !helpers.ContainsModel(teamTypesInterface, teamType1.Name) {
		t.Errorf("Expected team types to contain %v", teamType1.Name)
	}
	if !helpers.ContainsModel(teamTypesInterface, teamType2.Name) {
		t.Errorf("Expected team types to contain %v", teamType2.Name)
	}

	// Delete the team types
	db.Database.Unscoped().Delete(&teamType1)
	db.Database.Unscoped().Delete(&teamType2)

}
