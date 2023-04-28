package models

import (
	"testing"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveDelete tests the deletion of a dive
func TestDiveDelete(t *testing.T) {
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

	// Create the dive
	dive, err := models.DiveCreate(name, number, difficulty, uint64(divetypeID), uint64(divegroupID), uint64(boardtypeID), uint64(boardheightID))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}

	//Use the test data to create a new dive
	name = "Test Dive 2"
	number = 155
	difficulty = float32(2.5)
	divetypeID = dt2.ID
	divegroupID = dg2.ID
	boardtypeID = bt2.ID
	boardheightID = bh2.ID

	// Delete the dive
	err = models.DiveDelete(uint64(dive.ID))
	if err != nil {
		t.Fatalf("Error deleting dive: %v", err)
	}

	// Verify the dive was deleted
	_, err = models.DiveShow(uint64(dive.ID))
	if err == nil {
		t.Errorf("Dive with ID %d not deleted", dive.ID)
	}

	// Delete the dive
	db.Database.Unscoped().Delete(&dive)

	// Delete the test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
