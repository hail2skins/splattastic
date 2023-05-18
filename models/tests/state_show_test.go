package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestStateShow
func TestStateShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a state
	state, err := models.StateCreate("Test State", "TS")
	if err != nil {
		t.Fatal(err)
	}

	// Get the state
	state, err = models.StateShow(uint64(state.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the state name
	if state.Name != "Test State" {
		t.Errorf("handler returned unexpected state name: got %v want %v", state.Name, "TestState")
	}
	// Check the state code
	if state.Code != "TS" {
		t.Errorf("handler returned unexpected state code: got %v want %v", state.Code, "This is a short test description.")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&state)

}
