package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/hail2skins/splattastic/models/helpers"
)

// TestStatesGet tests the StatesGet function
func TestStatesGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two states
	state1, _ := models.StateCreate("Test State", "TS")
	state2, _ := models.StateCreate("Test Again", "TA")

	// Get all states
	states, err := models.StatesGet()
	if err != nil {
		t.Errorf("Error getting states: %v", err)
	}

	// Convert the states slice to a slice of interface{}
	statesInterface := make([]interface{}, len(states))
	for i, s := range states {
		statesInterface[i] = s
	}

	// Check that the states are the ones we created using containsModel function
	if !helpers.ContainsModel(statesInterface, state1.Name) {
		t.Errorf("Expected states to contain %v", state1.Name)
	}
	if !helpers.ContainsModel(statesInterface, state2.Name) {
		t.Errorf("Expected states to contain %v", state2.Name)
	}

	// Check that the states are the ones we created for code using containsModel function
	if !helpers.ContainsModelCode(statesInterface, state1.Code) {
		t.Errorf("Expected states to contain %v", state1.Code)
	}
	if !helpers.ContainsModelCode(statesInterface, state2.Code) {
		t.Errorf("Expected states to contain %v", state2.Code)
	}

	// Delete the states
	db.Database.Unscoped().Delete(&state1)
	db.Database.Unscoped().Delete(&state2)

}
