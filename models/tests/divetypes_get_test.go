package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/hail2skins/splattastic/models/helpers"
)

// TestDiveTypesGet tests the DiveTypesGet function
func TestDiveTypesGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two dive types
	diveType1, _ := models.DiveTypeCreate("TestDiveType1", "Q")
	diveType2, _ := models.DiveTypeCreate("TestDiveType2", "R")

	// Get all dive types
	diveTypes, err := models.DiveTypesGet()
	if err != nil {
		t.Errorf("Error getting dive types: %v", err)
	}

	// Convert the diveTypes slice to a slice of interface{}
	diveTypesInterface := make([]interface{}, len(diveTypes))
	for i, dt := range diveTypes {
		diveTypesInterface[i] = dt
	}

	// Check that the dive types are the ones we created using containsModel function
	if !helpers.ContainsModel(diveTypesInterface, diveType1.Name) {
		t.Errorf("Expected dive types to contain %v", diveType1.Name)
	}
	if !helpers.ContainsModel(diveTypesInterface, diveType2.Name) {
		t.Errorf("Expected dive types to contain %v", diveType2.Name)
	}

	// Delete the dive types
	db.Database.Unscoped().Delete(&diveType1)
	db.Database.Unscoped().Delete(&diveType2)
}
