package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/hail2skins/splattastic/models/helpers"
)

// TestMarkersGet tests the MarkersGet function
func TestMarkersGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two markers
	marker1, _ := models.MarkerCreate("TestMarker1", "This is a short test description.")
	marker2, _ := models.MarkerCreate("TestMarker2", "This is a short test description.")

	// Get all markers
	markers, err := models.MarkersGet()
	if err != nil {
		t.Errorf("Error getting markers: %v", err)
	}

	// Convert the markers slice to a slice of interface{}
	markersInterface := make([]interface{}, len(markers))
	for i, m := range markers {
		markersInterface[i] = m
	}

	// Check that the markers are the ones we created using containsModel function
	if !helpers.ContainsModel(markersInterface, marker1.Name) {
		t.Errorf("Expected markers to contain %v", marker1.Name)
	}
	if !helpers.ContainsModel(markersInterface, marker2.Name) {
		t.Errorf("Expected markers to contain %v", marker2.Name)
	}

	// Check that the markers are the ones we created for description using containsModel function
	if !helpers.ContainsModelDescription(markersInterface, marker1.Description) {
		t.Errorf("Expected markers to contain %v", marker1.Description)
	}
	if !helpers.ContainsModelDescription(markersInterface, marker2.Description) {
		t.Errorf("Expected markers to contain %v", marker2.Description)
	}

	// Delete the markers
	db.Database.Unscoped().Delete(&marker1)
	db.Database.Unscoped().Delete(&marker2)

}
