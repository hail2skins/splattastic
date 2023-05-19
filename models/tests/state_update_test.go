package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestStateUpdate is a test for the Update method of the State model
func TestStateUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new state
	state, err := models.StateCreate("TestState", "TS")
	if err != nil {
		t.Fatal(err)
	}

	// Update the state name
	newName := "UpdatedTestState"
	newCode := "UTS"
	err = state.Update(newName, newCode)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated state
	updatedState, err := models.StateShow(uint64(state.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedState.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedState.Name)
	}
	// Check if the code was updated
	if updatedState.Code != newCode {
		t.Errorf("expected updated code %v, got %v", newCode, updatedState.Code)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedState)
}
