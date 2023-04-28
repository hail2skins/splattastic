package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveGroupDelete tests the DiveTypeDelete function
func TestDiveGroupDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new dive group
	diveGroup, err := models.DiveGroupCreate("TestDiveGroup")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the dive group
	err = models.DiveGroupDelete(uint64(diveGroup.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Verify the dive group was deleted
	_, err = models.DiveGroupShow(uint64(diveGroup.ID))
	if err == nil {
		t.Errorf("Dive group with ID %d not deleted", diveGroup.ID)
	}

	// Delete the dive group
	db.Database.Unscoped().Delete(&diveGroup)
}
