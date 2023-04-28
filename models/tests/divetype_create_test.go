package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveTypeCreate tests the DiveTypeCreate function
func TestDiveTypeCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a board type
	diveType, err := models.DiveTypeCreate("TestDiveType", "Q")
	if err != nil {
		t.Errorf("DiveTypeCreate returned an error: %v", err)
	}

	// Test that the dive type was created
	if diveType.Name != "TestDiveType" || diveType.Letter != "Q" {
		t.Errorf("DiveTypeCreate returned an incorrect name: got %v want %v", diveType.Name, "TestDiveType")
		t.Errorf("DiveTypeCreate returned an incorrect letter: got %v want %v", diveType.Letter, "Q")
	}

	// Test duplicate dive type name
	dupDiveType, err := models.DiveTypeCreate("TestDiveType", "Q")
	if err == nil {
		t.Errorf("Expected error when creating dive type with duplicate name")
	}
	if dupDiveType != nil {
		t.Errorf("Expected dupDiveType to be nil")
	}

	// Test blank dive type name
	blankDiveType, err := models.DiveTypeCreate("", "")
	if err == nil {
		t.Errorf("Expected error when creating dive type with blank name")
	}
	if blankDiveType != nil {
		t.Errorf("Expected blankDiveType to be nil")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveType)
	db.Database.Unscoped().Delete(&dupDiveType)
	db.Database.Unscoped().Delete(&blankDiveType)
}
