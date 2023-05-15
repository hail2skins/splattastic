package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestMarkerDelete tests the MarkerDelete function
func TestMarkerDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new marker
	marker, err := models.MarkerCreate("TestMarker", "This is a short test description.")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the marker
	err = models.MarkerDelete(uint64(marker.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Verify the marker was deleted
	_, err = models.MarkerShow(uint64(marker.ID))
	if err == nil {
		t.Errorf("Marker with ID %d not deleted", marker.ID)
	}

	// Delete the marker
	db.Database.Unscoped().Delete(&marker)
}
