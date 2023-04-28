package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveTypeShow tests the DiveTypeShow function
func TestDiveTypeShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a dive type
	diveType, err := models.DiveTypeCreate("TestDiveType", "Q")
	if err != nil {
		t.Fatal(err)
	}

	// Get the dive type
	diveType, err = models.DiveTypeShow(uint64(diveType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the dive type name
	if diveType.Name != "TestDiveType" || diveType.Letter != "Q" {
		t.Errorf("handler returned unexpected dive type name: got %v want %v", diveType.Name, "TestDiveType")
		t.Errorf("handler returned unexpected dive type letter: got %v want %v", diveType.Letter, "Q")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveType)
}
