package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestBoardTypeUpdate is a test for the BoardTypeUpdate controller
func TestBoardTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new board type
	boardType, err := models.BoardTypeCreate("TestBoardType")
	if err != nil {
		t.Fatal(err)
	}

	// Update the board type name
	newName := "UpdatedTestBoardType"
	err = boardType.Update(newName)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated board type
	updatedBoardType, err := models.BoardTypeShow(uint64(boardType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	if updatedBoardType.Name != newName {
		t.Errorf("expected updated name %v, got %v", newName, updatedBoardType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedBoardType)
}
