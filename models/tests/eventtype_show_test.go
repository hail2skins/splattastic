package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventTypeShow tests the EventTypeShow function
func TestEventTypeShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create an event type
	eventType, err := models.EventTypeCreate("TestEventType")
	if err != nil {
		t.Fatal(err)
	}

	// Get the event type
	eventType, err = models.EventTypeShow(uint64(eventType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the event type name
	if eventType.Name != "TestEventType" {
		t.Errorf("handler returned unexpected event type name: got %v want %v", eventType.Name, "TestEventType")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&eventType)

}
