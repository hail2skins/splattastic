package models

import (
	"testing"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveUpdate tests the dive update functionality
func TestDiveUpdate(t *testing.T) {
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

	dive1, err := models.DiveCreate(name, number, difficulty, uint64(divetypeID), uint64(divegroupID), uint64(boardtypeID), uint64(boardheightID))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}

	// Update the dive with new values
	newName := "Updated Dive"
	newNumber := 254
	newDifficulty := float32(3.5)
	newDiveTypeID := dt2.ID
	newDiveGroupID := dg2.ID
	newBoardTypeID := bt2.ID
	newBoardHeightID := bh2.ID

	err = dive1.Update(newName, newNumber, newDifficulty, uint64(newDiveTypeID), uint64(newDiveGroupID), uint64(newBoardTypeID), uint64(newBoardHeightID))
	if err != nil {
		t.Fatalf("Error updating dive: %v", err)
	}

	// Fetch the updated dive from the database
	updatedDive, err := models.DiveShow(uint64(dive1.ID))
	if err != nil {
		t.Fatalf("Error fetching updated dive: %v", err)
	}

	// Check if the dive has been updated
	if updatedDive.Name != newName || updatedDive.Number != newNumber || updatedDive.Difficulty != newDifficulty ||
		updatedDive.DiveTypeID != uint64(newDiveTypeID) || updatedDive.DiveGroupID != uint64(newDiveGroupID) ||
		updatedDive.BoardTypeID != uint64(newBoardTypeID) || updatedDive.BoardHeightID != uint64(newBoardHeightID) {
		t.Errorf("Dive not updated correctly")
	}

	// Clean up the test dive and test data
	db.Database.Unscoped().Delete(dive1)
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}
