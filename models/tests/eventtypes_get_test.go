package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/hail2skins/splattastic/models/helpers"
)

// TestEventTypesGet tests the EventTypesGet function
func TestEventTypesGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two event types
	eventType1, _ := models.EventTypeCreate("TestEventType1")
	eventType2, _ := models.EventTypeCreate("TestEventType2")

	// Get all event types
	eventTypes, err := models.EventTypesGet()
	if err != nil {
		t.Errorf("Error getting event types: %v", err)
	}

	// Convert the eventTypes slice to a slice of interface{}
	eventTypesInterface := make([]interface{}, len(eventTypes))
	for i, et := range eventTypes {
		eventTypesInterface[i] = et
	}

	// Check that the event types are the ones we created using containsModel function
	if !helpers.ContainsModel(eventTypesInterface, eventType1.Name) {
		t.Errorf("Expected event types to contain %v", eventType1.Name)
	}
	if !helpers.ContainsModel(eventTypesInterface, eventType2.Name) {
		t.Errorf("Expected event types to contain %v", eventType2.Name)
	}

	// Delete the event types
	db.Database.Unscoped().Delete(&eventType1)
	db.Database.Unscoped().Delete(&eventType2)

}
