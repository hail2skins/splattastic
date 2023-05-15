package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestMarkerUpdate is a test for the MarkerUpdate controller
func TestMarkerUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new marker
	marker, err := models.MarkerCreate("TestMarker", "This is a short test description.")
	if err != nil {
		t.Fatal(err)
	}

	// Update the marker name
	newName := "UpdatedTestMarker"
	newDescription := "This is an updated description."
	err = marker.Update(newName, newDescription)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated marker
	updatedMarker, err := models.MarkerShow(uint64(marker.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedMarker.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedMarker.Name)
	}
	// Check if the description was updated
	if updatedMarker.Description != newDescription {
		t.Errorf("expected updated description %v, got %v", newDescription, updatedMarker.Description)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedMarker)
}
