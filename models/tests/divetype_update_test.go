package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveTypeUpdate is a test for the DiveTypeUpdate controller
func TestDiveTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new dive type
	diveType, err := models.DiveTypeCreate("TestDiveType")
	if err != nil {
		t.Fatal(err)
	}

	// Update the dive type name
	newName := "UpdatedTestDiveType"
	err = diveType.Update(newName)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated dive type
	updatedDiveType, err := models.DiveTypeShow(uint64(diveType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedDiveType.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedDiveType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedDiveType)
}
