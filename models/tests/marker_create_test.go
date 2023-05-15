package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestMarkerCreate tests the MarkerCreate function
func TestMarkerCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a marker
	marker, err := models.MarkerCreate("TestMarker")
	if err != nil {
		t.Errorf("Error creating marker: %v", err)
	}

	// Test that the marker was created
	if marker.Name != "TestMarker" {
		t.Errorf("MarkerCreate returned an incorrect name: got %v want %v", marker.Name, "TestMarker")
	}

	// Test duplicate marker name
	dupMarker, err := models.MarkerCreate("TestMarker")
	if err == nil {
		t.Errorf("Expected error when creating marker with duplicate name")
	}

	// Test blank marker name
	blankMarker, err := models.MarkerCreate("")
	if err == nil {
		t.Errorf("Expected error when creating marker with blank name")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&marker)
	db.Database.Unscoped().Delete(&dupMarker)
	db.Database.Unscoped().Delete(&blankMarker)

}
