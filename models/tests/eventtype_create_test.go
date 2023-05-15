package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventTypeCreate tests the EventTypeCreate function
func TestEventTypeCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create an event type
	eventType, err := models.EventTypeCreate("TestEventType")
	if err != nil {
		t.Errorf("Error creating event type: %v", err)
	}

	// Test that the event type was created
	if eventType.Name != "TestEventType" {
		t.Errorf("EventTypeCreate returned an incorrect name: got %v want %v", eventType.Name, "TestEventType")
	}

	// Test duplicate event type name
	dupEventType, err := models.EventTypeCreate("TestEventType")
	if err == nil {
		t.Errorf("Expected error when creating event type with duplicate name")
	}

	// Test blank event type name
	blankEventType, err := models.EventTypeCreate("")
	if err == nil {
		t.Errorf("Expected error when creating event type with blank name")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&eventType)
	db.Database.Unscoped().Delete(&dupEventType)
	db.Database.Unscoped().Delete(&blankEventType)

}
