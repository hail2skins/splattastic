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
	boardTypesInterface := make([]interface{}, len(boardTypes))
	for i, bt := range boardTypes {
		boardTypesInterface[i] = bt
	}

	// Delete the board types
	db.Database.Unscoped().Delete(&boardType1)
	db.Database.Unscoped().Delete(&boardType2)
}
