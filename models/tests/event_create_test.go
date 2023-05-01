package models

import (
	"testing"
	"time"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventCreate tests the EventCreate function
// Need to create the associated records first
func TestEventCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// // Insert test data and defer cleanup
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
	date := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	against := "Test Against"
	userID := user.ID
	eventtypeID := et1.ID

	// Test with 0 dives
	event, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), []uint64{})
	if err != nil {
		t.Fatalf("Error creating event with 0 dives: %v", err)
	}

	// Check if event was created and associated with the user
	if event.Name != name ||
		event.Location != location ||
		event.Date != date ||
		event.Against != against ||
		event.UserID != uint64(userID) ||
		event.EventTypeID != uint64(eventtypeID) {
		t.Fatal("Event was not created correctly")
	}

	// Clean up event
	db.Database.Unscoped().Delete(&event)

	// Test with 1 dive
	event1, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), []uint64{uint64(dive.ID)})
	if err != nil {
		t.Fatalf("Error creating event with 1 dive: %v", err)
	}

	// Check if the event is associated with the dive
	if !event1.HasDive(dive) {
		t.Fatal("Event should be associated with the dive")
	}

	// Clean up event1
	db.Database.Unscoped().Delete(&event1)
	// Clean up UserEventDive
	deleteUserEventDives()

	// Test with 2 dives
	event2, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), []uint64{uint64(dive.ID), uint64(dive2.ID)})
	if err != nil {
		t.Fatalf("Error creating event with 2 dives: %v", err)
	}

	// Check if the event is associated with both dives
	if !event2.HasDive(dive) || !event2.HasDive(dive2) {
		t.Fatal("Event should be associated with both dives")
	}
	// Clean up event2
	db.Database.Unscoped().Delete(&event2)

	// Delete user event dives
	deleteUserEventDives()

	// Delete dives
	db.Database.Unscoped().Delete(&dive)
	db.Database.Unscoped().Delete(&dive2)

	// Delete event
	db.Database.Unscoped().Delete(&event)

	// Delete user
	db.Database.Unscoped().Delete(&user)

	// Delete user type
	db.Database.Unscoped().Delete(&ut1)

	// Delete event type
	db.Database.Unscoped().Delete(&et1)

	// Remove test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}

func deleteUserEventDives() error {
	var userEventDives []models.UserEventDive
	result := db.Database.Find(&userEventDives)
	if result.Error != nil {
		return result.Error
	}

	//fmt.Println("UserEventDives before deletion:", userEventDives)

	for _, userEventDive := range userEventDives {
		db.Database.Unscoped().Delete(&userEventDive)
	}

	// Check if the user event dives are deleted
	var userEventDivesAfterDeletion []models.UserEventDive
	result = db.Database.Find(&userEventDivesAfterDeletion)
	if result.Error != nil {
		return result.Error
	}
	//fmt.Println("UserEventDives after deletion:", userEventDivesAfterDeletion)

	return nil
}
