package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestStateDelete tests the StateDelete function
func TestStateDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new state
	state, err := models.StateCreate("TestState", "TS")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the state
	err = models.StateDelete(uint64(state.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Verify the state was deleted
	_, err = models.StateShow(uint64(state.ID))
	if err == nil {
		t.Errorf("State with ID %d not deleted", state.ID)
	}

	// Delete the state
	db.Database.Unscoped().Delete(&state)
}
