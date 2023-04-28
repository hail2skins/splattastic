package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveTypeDelete tests the DiveTypeDelete function
func TestDiveTypeDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new dive type
	diveType, err := models.DiveTypeCreate("TestDiveType", "Q")
	if err != nil {
		t.Fatal(err)
	}

	// Verify the dive type was created
	_, err = models.DiveTypeShow(uint64(diveType.ID))
	if err != nil {
		t.Errorf("Dive type with ID %d not created", diveType.ID)
	}

	// Delete the dive type
	db.Database.Unscoped().Delete(&diveType)

}
