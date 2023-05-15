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
	marker, err := models.MarkerCreate("TestMarker")
	if err != nil {
		t.Fatal(err)
	}

	// Update the marker name
	newName := "UpdatedTestMarker"
	err = marker.Update(newName)
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

	// Cleanup
	db.Database.Unscoped().Delete(&updatedMarker)
}
