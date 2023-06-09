package models

import (
	"testing"
	"time"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventDelete tests the event delete function
func TestEventDelete(t *testing.T) {
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

	// Create an Event Type
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

	// Test with 1 dives
	event, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), []uint64{uint64(dive.ID)})
	if err != nil {
		t.Fatalf("Error creating event with 0 dives: %v", err)
	}

	// Delete the event
	err = models.EventDelete(uint64(event.ID))
	if err != nil {
		t.Fatalf("Error deleting event: %v", err)
	}

	// Verify the event is deleted
	_, err = models.EventShow(uint64(event.ID))
	if err == nil {
		t.Fatalf("Event not deleted")
	}

	// Verify the User Event Dive is deleted
	userEventDives, _ := models.GetDivesForEvent(uint64(event.ID))

	for _, userEventDive := range userEventDives {
		if userEventDive.DeletedAt.Valid {
			t.Logf("DeletedAt: %v", userEventDive.DeletedAt.Time)
		} else {
			t.Fatalf("User Event Dive not deleted: %v", userEventDive)
		}
	}

	// Delete the User Event Dives
	_ = models.DeleteUserEventDivesByEventID(event.ID)

	// Delete the event
	db.Database.Unscoped().Delete(&event)
	// Delete the event type
	db.Database.Unscoped().Delete(&et1)
	// Delete the dives
	db.Database.Unscoped().Delete(&dive)
	db.Database.Unscoped().Delete(&dive2)
	// Delete the user
	db.Database.Unscoped().Delete(&user)
	// Delete the user type
	db.Database.Unscoped().Delete(&ut1)
	// Delete the test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
