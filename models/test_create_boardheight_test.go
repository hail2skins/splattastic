package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

func TestCreateBoardHeight(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Test data
	testHeight := float32(9.5)

	// Call the CreateBoardHeight function
	boardHeight, err := CreateBoardHeight(testHeight)

	// Assert there was no error and a board height was returned
	assert.NoError(t, err)
	assert.NotNil(t, boardHeight)

	// Assert that the returned board height has the correct height value
	assert.Equal(t, testHeight, boardHeight.Height)

	// Query the database to check if the board height was saved correctly
	var dbBoardHeight BoardHeight
	err = db.Database.First(&dbBoardHeight, boardHeight.ID).Error

	// Assert there was no error and the board height was found in the database
	assert.NoError(t, err)
	assert.NotNil(t, &dbBoardHeight)

	// Assert that the database board height has the correct height value
	assert.Equal(t, testHeight, dbBoardHeight.Height)

	// Cleanup
	db.Database.Unscoped().Delete(boardHeight)
}
