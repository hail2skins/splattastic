package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveGroupUpdate is a test for the DiveGroupUpdate controller
func TestDiveGroupUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new dive group
	diveGroup, err := models.DiveGroupCreate("TestDiveGroup")
	if err != nil {
		t.Fatal(err)
	}

	// Update the dive group name
	newName := "UpdatedTestDiveGroup"
	err = diveGroup.Update(newName)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated dive group
	updatedDiveGroup, err := models.DiveGroupShow(uint64(diveGroup.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedDiveGroup.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedDiveGroup.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedDiveGroup)
}
