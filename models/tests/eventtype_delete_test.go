package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventTypeDelete tests the EventTypeDelete function
func TestEventTypeDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new event type
	eventType, err := models.EventTypeCreate("TestEventType")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the event type
	err = models.EventTypeDelete(uint64(eventType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Verify the event type was deleted
	_, err = models.EventTypeShow(uint64(eventType.ID))
	if err == nil {
		t.Errorf("Event type with ID %d not deleted", eventType.ID)
	}

	// Delete the event type
	db.Database.Unscoped().Delete(&eventType)
}
