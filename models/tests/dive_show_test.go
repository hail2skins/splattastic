package models

import (
	"testing"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveShow tests the DiveShow function
func TestDiveShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Use the test data to create a new dive
	name := "Test Dive"
	number := 154
	difficulty := float32(2.5)
	divetypeID := dt1.ID
	divegroupID := dg1.ID
	boardtypeID := bt1.ID
	boardheightID := bh1.ID

	dive, err := models.DiveCreate(name, number, difficulty, uint64(divetypeID), uint64(divegroupID), uint64(boardtypeID), uint64(boardheightID))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}

	// Get the dive
	dive, err = models.DiveShow(uint64(dive.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the dive name
	if dive.Name != name || dive.Number != number || dive.Difficulty != difficulty || dive.DiveTypeID != uint64(divetypeID) || dive.DiveGroupID != uint64(divegroupID) || dive.BoardTypeID != uint64(boardtypeID) || dive.BoardHeightID != uint64(boardheightID) {
		t.Error("DiveShow returned incorrect data")
	}

	// Delete the created dive
	// Hard-delete the created dive
	err = db.Database.Unscoped().Delete(dive).Error
	if err != nil {
		t.Fatalf("Error hard-deleting dive: %v", err)
	}

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
