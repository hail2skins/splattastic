package models

import (
	"testing"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDivesGet tests the DiveGet function
func TestDivesGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Use the test data to create a new dive
	name := "Test Dive"
	number := 101
	difficulty := float32(2.5)
	divetypeID := dt1.ID
	divegroupID := dg1.ID
	boardtypeID := bt1.ID
	boardheightID := bh1.ID

	// Use this test data to create a second dive
	name2 := "Test Dive 2"
	number2 := 102
	difficulty2 := float32(3.5)
	divetypeID2 := dt2.ID
	divegroupID2 := dg2.ID
	boardtypeID2 := bt2.ID
	boardheightID2 := bh2.ID

	// Create the first dive
	dive, err := models.DiveCreate(name, number, difficulty, uint64(divetypeID), uint64(divegroupID), uint64(boardtypeID), uint64(boardheightID))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}
	// Create the second dive
	dive2, err := models.DiveCreate(name2, number2, difficulty2, uint64(divetypeID2), uint64(divegroupID2), uint64(boardtypeID2), uint64(boardheightID2))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}

	// Test DivesGet
	dives, err := models.DivesGet()
	if err != nil {
		t.Fatalf("Error getting dives: %v", err)
	}

	// Check if created dives are in the result set and their related records are loaded correctly
	diveFound := false
	dive2Found := false

	for _, d := range dives {
		if d.ID == dive.ID {
			diveFound = true
			if uint64(d.DiveType.ID) != dive.DiveTypeID || uint64(d.DiveGroup.ID) != dive.DiveGroupID || uint64(d.BoardType.ID) != dive.BoardTypeID || uint64(d.BoardHeight.ID) != dive.BoardHeightID {
				t.Error("DivesGet did not preload related records correctly")
			}
		}
		if d.ID == dive2.ID {
			dive2Found = true
			if uint64(d.DiveType.ID) != dive2.DiveTypeID || uint64(d.DiveGroup.ID) != dive2.DiveGroupID || uint64(d.BoardType.ID) != dive2.BoardTypeID || uint64(d.BoardHeight.ID) != dive2.BoardHeightID {
				t.Error("DivesGet did not preload related records correctly")
			}
		}
	}

	if !diveFound || !dive2Found {
		t.Error("DivesGet did not return the created dives")
	}

	// Delete the created dive
	// Hard-delete the created dive
	err = db.Database.Unscoped().Delete(dive).Error
	if err != nil {
		t.Fatalf("Error hard-deleting dive: %v", err)
	}

	// Delete the created dive
	// Hard-delete the created dive
	err = db.Database.Unscoped().Delete(dive2).Error
	if err != nil {
		t.Fatalf("Error hard-deleting dive: %v", err)
	}

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
