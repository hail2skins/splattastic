package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestStateCreate tests the StateCreate function
func TestStateCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a state
	state, err := models.StateCreate("TestState", "TS")
	if err != nil {
		t.Errorf("Error creating state: %v", err)
	}

	// Test that the state was created
	if state.Name != "TestState" {
		t.Errorf("StateCreate returned an incorrect name: got %v want %v", state.Name, "TestState")
	}

	// Test that the state code is as expected
	if state.Code != "TS" {
		t.Errorf("StateCreate returned an incorrect code: got %v want %v", state.Code, "TS")
	}

	// Test duplicate state name
	dupState, err := models.StateCreate("TestState", "TS")
	if err == nil {
		t.Errorf("Expected error when creating state with duplicate name")
	}

	// Test blank state name
	blankState, err := models.StateCreate("", "TS")
	if err == nil {
		t.Errorf("Expected error when creating state with blank name")
	}

	// Test blank state code
	blankCode, err := models.StateCreate("TestState", "")
	if err == nil {
		t.Errorf("Expected error when creating state with blank code")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&state)
	db.Database.Unscoped().Delete(&dupState)
	db.Database.Unscoped().Delete(&blankState)
	db.Database.Unscoped().Delete(&blankCode)

}
