package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func TestBoardTypeDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a new board type
	boardType, err := models.BoardTypeCreate("TestBoardType")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the board type
	err = models.BoardTypeDelete(uint64(boardType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Verify the board type was deleted
	_, err = models.BoardTypeShow(uint64(boardType.ID))
	if err == nil {
		t.Errorf("Board type with ID %d not deleted", boardType.ID)
	}

	// Teardown
	db.Database.Unscoped().Delete(&boardType)
}
