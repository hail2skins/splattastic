package models

import (
	"fmt"
	"testing"
	"time"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestGetLastFiveEvents tests the GetLastFiveEvents function
func TestGetLastFiveEvents(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a test user type
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
	eventDate = eventDate.Add(time.Second)
	event1, _ := models.EventCreate("Test Event1", "Test Location1", eventDate, "Test Against1", uint64(user.ID), uint64(et1.ID), []uint64{})
	eventDate = eventDate.Add(time.Second)
	event2, _ := models.EventCreate("Test Event2", "Test Location2", eventDate, "Test Against2", uint64(user.ID), uint64(et1.ID), []uint64{})
	eventDate = eventDate.Add(time.Second)
	event3, _ := models.EventCreate("Test Event3", "Test Location3", eventDate, "Test Against3", uint64(user.ID), uint64(et1.ID), []uint64{})
	eventDate = eventDate.Add(time.Second)
	event4, _ := models.EventCreate("Test Event4", "Test Location4", eventDate, "Test Against4", uint64(user.ID), uint64(et1.ID), []uint64{})
	eventDate = eventDate.Add(time.Second)
	event5, _ := models.EventCreate("Test Event5", "Test Location5", eventDate, "Test Against5", uint64(user.ID), uint64(et1.ID), []uint64{})

	// Get the last five events for the user
	events, err := models.GetLastFiveEvents(uint64(user.ID))
	for _, e := range *events {
		fmt.Printf("Event: %+v\n", e)
	}

	if err != nil {
		t.Fatalf("Error getting user events: %v", err)
	}

	// Check if created events are in the slice
	eventFound := false
	event1Found := false
	event2Found := false
	event3Found := false
	event4Found := false
	event5Found := false

	// Loop through the events
	for _, e := range *events {
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
		if e.ID == event2.ID {
			event2Found = true
			if uint64(e.EventType.ID) != uint64(et1.ID) {
				t.Fatalf("Event type does not match")
			}
		}
		if e.ID == event3.ID {
			event3Found = true
			if uint64(e.EventType.ID) != uint64(et1.ID) {
				t.Fatalf("Event type does not match")
			}
		}
		if e.ID == event4.ID {
			event4Found = true
			if uint64(e.EventType.ID) != uint64(et1.ID) {
				t.Fatalf("Event type does not match")
			}
		}
		if e.ID == event5.ID {
			event5Found = true
			if uint64(e.EventType.ID) != uint64(et1.ID) {
				t.Fatalf("Event type does not match")
			}
		}
	}

	// Check if the five events we want are in the list

	if !event5Found {
		t.Fatalf("Event5 not found")
	}
	if !event4Found {
		t.Fatalf("Event4 not found")
	}
	if !event3Found {
		t.Fatalf("Event3 not found")
	}
	if !event2Found {
		t.Fatalf("Event2 not found")
	}
	if !event1Found {
		t.Fatalf("Event1 not found")
	}

	// Check if the sixth event is not in the list
	if eventFound {
		t.Fatalf("Event found")
	}

	// Delete the events
	db.Database.Unscoped().Delete(&event)
	db.Database.Unscoped().Delete(&event1)
	db.Database.Unscoped().Delete(&event2)
	db.Database.Unscoped().Delete(&event3)
	db.Database.Unscoped().Delete(&event4)
	db.Database.Unscoped().Delete(&event5)

	// Delete the user
	db.Database.Unscoped().Delete(&user)

	// Delete the event type
	db.Database.Unscoped().Delete(&et1)

	// Delete the user type
	db.Database.Unscoped().Delete(&ut1)

}
