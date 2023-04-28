package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestBoardTypeShow tests the BoardTypeShow function
func TestBoardTypeShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a board type
	boardType, err := models.BoardTypeCreate("TestBoardType")
	if err != nil {
		t.Fatal(err)
	}

	// Get the board type
	boardType, err = models.BoardTypeShow(uint64(boardType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the board type name
	if boardType.Name != "TestBoardType" {
		t.Errorf("handler returned unexpected board type name: got %v want %v", boardType.Name, "TestBoardType")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&boardType)
}
