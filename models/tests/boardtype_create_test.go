package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestBoardTypeCreate tests the BoardTypeCreate function
func TestBoardTypeCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a board type
	boardType, err := models.BoardTypeCreate("Test Board Type")
	if err != nil {
		t.Errorf("Error creating board type: %v", err)
	}
	if boardType.Name != "Test Board Type" {
		t.Errorf("Expected board type name to be 'Test Board Type', got '%v'", boardType.Name)
	}

	// Test duplicate board type name
	dupBoardType, err := models.BoardTypeCreate("Test Board Type")
	if err == nil {
		t.Errorf("Expected error when creating board type with duplicate name")
	}
	if dupBoardType != nil {
		t.Errorf("Expected dupBoardType to be nil")
	}

	// Test blank board type name
	blankBoardType, err := models.BoardTypeCreate("")
	if err == nil {
		t.Errorf("Expected error when creating board type with blank name")
	}
	if blankBoardType != nil {
		t.Errorf("Expected blankBoardType to be nil")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&boardType)
	db.Database.Unscoped().Delete(&dupBoardType)
	db.Database.Unscoped().Delete(&blankBoardType)
}
