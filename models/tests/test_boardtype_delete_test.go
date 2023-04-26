package models

import (
	"os"
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func TestBoardTypeDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

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
