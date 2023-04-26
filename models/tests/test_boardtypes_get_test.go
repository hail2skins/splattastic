package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestBoardTypesGet tests the BoardTypesGet function
func TestBoardTypesGet(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create two board types
	boardType1, _ := models.BoardTypeCreate("TestBoardType1")
	boardType2, _ := models.BoardTypeCreate("TestBoardType2")

	// Get all board types
	boardTypes, err := models.BoardTypesGet()
	if err != nil {
		t.Errorf("Error getting board types: %v", err)
	}

	// Check that the board types are the ones we created using containsBoardType function
	if !containsBoardType(boardTypes, boardType1.Name) {
		t.Errorf("Expected board types to contain %v", boardType1.Name)
	}
	if !containsBoardType(boardTypes, boardType2.Name) {
		t.Errorf("Expected board types to contain %v", boardType2.Name)
	}

	// Delete the board types
	db.Database.Unscoped().Delete(&boardType1)
	db.Database.Unscoped().Delete(&boardType2)
}

// containsBoardType function to ensure a board type is in a slice of board types regardless of order
func containsBoardType(boardTypes []models.BoardType, name string) bool {
	for _, b := range boardTypes {
		if b.Name == name {
			return true
		}
	}
	return false
}
