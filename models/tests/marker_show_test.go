package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestMarkerShow tests the MarkerShow function
func TestMarkerShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a marker
	marker, err := models.MarkerCreate("TestMarker")
	if err != nil {
		t.Fatal(err)
	}

	// Get the marker
	marker, err = models.MarkerShow(uint64(marker.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the marker name
	if marker.Name != "TestMarker" {
		t.Errorf("handler returned unexpected marker name: got %v want %v", marker.Name, "TestMarker")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&marker)

}
