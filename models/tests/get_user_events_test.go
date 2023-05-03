package models

import (
	"testing"
	"time"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestGetUserEvents retrieves all events for a user.
// This will NOT populate UserEventDives, so not testing for that.
func TestGetUserEvents(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Insert test data and defer cleanup
	// Not really necessary here but leaving it in for now
	//dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two test dives (similar to what you did in the model test)
	// Not really necessary here but leaving it in for now
	//dive1, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))
	//dive2, _ := models.DiveCreate("Test Dive 2", 155, 3.5, uint64(dt2.ID), uint64(dg2.ID), uint64(bt2.ID), uint64(bh2.ID))

	// Create a UserType
	ut1, err := models.CreateUserType("Test User Type")
	if err != nil {
		t.Fatalf("Error creating user type: %v", err)
	}

	// Create a User
	user, err := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	// Create a test event type
	et1, err := models.EventTypeCreate("Test Event Type")
	if err != nil {
		t.Fatalf("Error creating event type: %v", err)
	}

	// Create multiple events for the user
	eventDate := time.Now()
	event, _ := models.EventCreate("Test Event", "Test Location", eventDate, "Test Against", uint64(user.ID), uint64(et1.ID), []uint64{})
	event1, _ := models.EventCreate("Test Event1", "Test Location1", eventDate, "Test Against1", uint64(user.ID), uint64(et1.ID), []uint64{})

	// Get all events for the user
	events, err := models.GetUserEvents(uint64(user.ID))
	if err != nil {
		t.Fatalf("Error getting user events: %v", err)
	}

	// Check if created events are in the slice
	eventFound := false
	event1Found := false

	// Loop through the events
	for _, e := range events {
		// Check if the event is one of the created events
		if e.ID == event.ID {
			eventFound = true
			if uint64(e.EventType.ID) != uint64(et1.ID) {
				t.Fatalf("Event type does not match")
			}
		}
		if e.ID == event1.ID {
			event1Found = true
			if uint64(e.EventType.ID) != uint64(et1.ID) {
				t.Fatalf("Event type does not match")
			}
		}
	}

	// Check if the events were found
	if !eventFound || !event1Found {
		t.Fatalf("Events not found")
	}

	// Delete the events
	db.Database.Unscoped().Delete(&event)
	db.Database.Unscoped().Delete(&event1)

	// Delete the event type
	db.Database.Unscoped().Delete(&et1)

	// Delete the dives
	//db.Database.Unscoped().Delete(&dive1)
	//db.Database.Unscoped().Delete(&dive2)

	// Delete the user
	db.Database.Unscoped().Delete(&user)

	// Delete the user type
	db.Database.Unscoped().Delete(&ut1)

	// Clean the test data
	//helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
