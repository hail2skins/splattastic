package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveGroupCreate tests the DiveTypeCreate function
func TestDiveGroupCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a board type
	diveGroup, err := models.DiveGroupCreate("TestDiveGroup")
	if err != nil {
		t.Errorf("DiveGroupCreate returned an error: %v", err)
	}

	// Test that the dive type was created
	if diveGroup.Name != "TestDiveGroup" {
		t.Errorf("DiveGroupCreate returned an incorrect name: got %v want %v", diveGroup.Name, "Test Dive Group")
	}

	// Test duplicate dive type name
	dupDiveGroup, err := models.DiveGroupCreate("TestDiveGroup")
	if err == nil {
		t.Errorf("Expected error when creating dive group with duplicate name")
	}

	// Test blank dive type name
	blankDiveGroup, err := models.DiveGroupCreate("")
	if err == nil {
		t.Errorf("Expected error when creating dive group with blank name")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveGroup)
	db.Database.Unscoped().Delete(&dupDiveGroup)
	db.Database.Unscoped().Delete(&blankDiveGroup)

}
