package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

// TestBoardHeightShow is a function that returns a single board height
func TestBoardHeightShow(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a board height
	testBoardHeight := BoardHeight{Height: 32}
	db.Database.Create(&testBoardHeight)

	// Test BoardHeightShow with valid ID
	boardHeight, err := BoardHeightShow(uint64(testBoardHeight.ID))
	assert.NoError(t, err)
	assert.NotNil(t, boardHeight)
	assert.Equal(t, testBoardHeight.ID, boardHeight.ID)
	assert.Equal(t, testBoardHeight.Height, boardHeight.Height)

	// Test BoardHeightShow with invalid ID
	_, err = BoardHeightShow(0) // ID 0 is invalid
	assert.Error(t, err)

	// Clean up
	db.Database.Unscoped().Delete(&testBoardHeight)
}
