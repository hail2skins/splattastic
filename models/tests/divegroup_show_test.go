package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveGroupShow tests the DiveGroupShow function
func TestDiveGroupShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a dive group
	diveGroup, err := models.DiveGroupCreate("TestDiveGroup")
	if err != nil {
		t.Fatal(err)
	}

	// Get the dive group
	diveGroup, err = models.DiveGroupShow(uint64(diveGroup.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the dive group name
	if diveGroup.Name != "TestDiveGroup" {
		t.Errorf("handler returned unexpected dive group name: got %v want %v", diveGroup.Name, "TestDiveGroup")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveGroup)
}
