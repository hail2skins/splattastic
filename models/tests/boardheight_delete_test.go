package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestBoardHeightDelete is a model for testing the deletion of a board height
func TestBoardHeightDelete(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create a test board height
	testBoardHeight := models.BoardHeight{Height: 20}
	db.Database.Create(&testBoardHeight)

	assert.NotNil(t, testBoardHeight)

	// Soft delete the board height
	err := models.BoardHeightDelete(uint64(testBoardHeight.ID))
	assert.NoError(t, err)

	// Try to get the board height
	var softDeletedBoardHeight models.BoardHeight
	result := db.Database.Unscoped().First(&softDeletedBoardHeight, testBoardHeight.ID)
	assert.NoError(t, result.Error) // No error is expected when fetching the soft-deleted record
	assert.NotNil(t, softDeletedBoardHeight.DeletedAt)

	// Cleanup
	db.Database.Unscoped().Delete(&testBoardHeight)

}
