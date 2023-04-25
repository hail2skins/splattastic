package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

// TestUpdateBoardHeight
func TestUpdateBoardHeight(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a board height
	testBoardHeight := BoardHeight{Height: 18}
	db.Database.Create(&testBoardHeight)

	// Update the board height
	newHeight := float32(20)
	err := testBoardHeight.Update(newHeight)
	assert.NoError(t, err)

	// Get the board height
	var updatedBoardHeight BoardHeight
	db.Database.First(&updatedBoardHeight, testBoardHeight.ID)

	// Verify the board height was updated
	assert.Equal(t, newHeight, updatedBoardHeight.Height)

	// Delete the board height
	db.Database.Unscoped().Delete(&testBoardHeight)
}
