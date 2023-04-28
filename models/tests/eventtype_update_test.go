package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventTypeUpdate is a test for the EventTypeUpdate controller
func TestEventTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new event type
	eventType, err := models.EventTypeCreate("TestEventType")
	if err != nil {
		t.Fatal(err)
	}

	// Update the event type name
	newName := "UpdatedTestEventType"
	err = eventType.Update(newName)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated event type
	updatedEventType, err := models.EventTypeShow(uint64(eventType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedEventType.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedEventType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedEventType)
}
