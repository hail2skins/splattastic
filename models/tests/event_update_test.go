package models

import (
	"testing"
	"time"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventUpdate tests the EventUpdate function
// Need to create the associated records first
func TestEventUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two dives
	// Use the test data to create a new dive
	name := "Test Dive"
	number := 154
	difficulty := float32(2.5)
	divetypeID := dt1.ID
	divegroupID := dg1.ID
	boardtypeID := bt1.ID
	boardheightID := bh1.ID

	// Use this test data to create a second dive
	name2 := "Test Dive 2"
	number2 := 155
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

	// Create a UserType
	ut1, err := models.CreateUserType("Test User Type")
	if err != nil {
		t.Fatalf("Error creating user type: %v", err)
	}

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)

	// Create EventType
	et1, err := models.EventTypeCreate("Test Event Type")
	if err != nil {
		t.Fatalf("Error creating event type: %v", err)
	}

	// Use the test data to create a new event
	name = "Test Event"
	location := "Test Location"
	date := time.Date(2018, 1, 2, 0, 0, 0, 0, time.UTC)
	against := "Test Against"
	eventtypeID := et1.ID
	userID := user.ID

	testCases := []struct {
		name         string
		initialDives []uint64
		newDives     []uint64
	}{
		{
			name:         "Update event with 0 dives to 1",
			initialDives: []uint64{},
			newDives:     []uint64{uint64(dive.ID)},
		},
		{
			name:         "Update event with 0 dives to 2",
			initialDives: []uint64{},
			newDives:     []uint64{uint64(dive.ID), uint64(dive2.ID)},
		},
		{
			name:         "Update event with 1 dive, removing it",
			initialDives: []uint64{uint64(dive.ID)},
			newDives:     []uint64{},
		},
		{
			name:         "Update event with 2 dives, removing one",
			initialDives: []uint64{uint64(dive.ID), uint64(dive2.ID)},
			newDives:     []uint64{uint64(dive.ID)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new event with the initial dives
			event, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), tc.initialDives)
			if err != nil {
				t.Fatalf("Error creating event: %v", err)
			}

			// Update the event with the new dives
			updatedEvent, err := event.Update(name, location, date, against, uint64(eventtypeID), tc.newDives)
			if err != nil {
				t.Fatalf("Error updating event: %v", err)
			}

			// Get the updated dives for the event
			updatedDives, err := models.GetDivesForEvent(uint64(updatedEvent.ID))
			if err != nil {
				t.Fatalf("Error retrieving updated dives for event: %v", err)
			}

			// Check if the updated dives match the expected new dives
			if len(updatedDives) != len(tc.newDives) {
				t.Fatalf("Updated dives count doesn't match expected count: got %v, want %v", len(updatedDives), len(tc.newDives))
			}
			for i, dive := range updatedDives {
				if uint64(dive.ID) != tc.newDives[i] {
					t.Errorf("Updated dive ID doesn't match expected dive ID: got %v, want %v", dive.ID, tc.newDives[i])
				}
			}

			// Find and delete UserEventDives created during EventCreate and event.Update
			var allUserEventDives []models.UserEventDive
			db.Database.Where("event_id = ?", updatedEvent.ID).Find(&allUserEventDives)
			for _, userEventDive := range allUserEventDives {
				db.Database.Unscoped().Delete(&userEventDive)
			}

			if err := models.DeleteUserEventDivesByEventID(event.ID); err != nil {
				t.Errorf("Error deleting user event dives: %v", err)
			}

			// Cleanup: delete the event
			db.Database.Unscoped().Delete(&event)
		})
	}

	// Cleanup
	// Delete user event dives
	helpers.DeleteUserEventDives()
	// Delete dives
	db.Database.Unscoped().Delete(&dive)
	db.Database.Unscoped().Delete(&dive2)
	// Delete user
	db.Database.Unscoped().Delete(&user)
	// Delete user type
	db.Database.Unscoped().Delete(&ut1)
	// Delete event type
	db.Database.Unscoped().Delete(&et1)
	// Remove test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
