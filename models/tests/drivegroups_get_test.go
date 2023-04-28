package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/hail2skins/splattastic/models/helpers"
)

// TestDiveGroupsGet tests the DiveGroupsGet function
func TestDiveGroupsGet(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create two dive groups
	diveGroup1, _ := models.DiveGroupCreate("TestDiveGroup1")
	diveGroup2, _ := models.DiveGroupCreate("TestDiveGroup2")

	// Get all dive groups
	diveGroups, err := models.DiveGroupsGet()
	if err != nil {
		t.Errorf("Error getting dive groups: %v", err)
	}

	// Convert the diveGroups slice to a slice of interface{}
	diveGroupsInterface := make([]interface{}, len(diveGroups))
	for i, dg := range diveGroups {
		diveGroupsInterface[i] = dg
	}

	// Check that the dive groups are the ones we created using containsModel function
	if !helpers.ContainsModel(diveGroupsInterface, diveGroup1.Name) {
		t.Errorf("Expected dive groups to contain %v", diveGroup1.Name)
	}
	if !helpers.ContainsModel(diveGroupsInterface, diveGroup2.Name) {
		t.Errorf("Expected dive groups to contain %v", diveGroup2.Name)
	}

	// Delete the dive groups
	db.Database.Unscoped().Delete(&diveGroup1)
	db.Database.Unscoped().Delete(&diveGroup2)

}
