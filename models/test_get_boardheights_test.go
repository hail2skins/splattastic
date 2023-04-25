package models

import (
	"testing"

	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

// TestGetBoardHeights is a function that tests the GetBoardHeights function
func TestGetBoardHeights(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Create user types
	boardHeight1 := BoardHeight{Height: 2}
	boardHeight2 := BoardHeight{Height: 4.5}
	db.Database.Create(&boardHeight1)
	db.Database.Create(&boardHeight2)

	// Call GetUserTypes function
	boardHeights, err := GetBoardHeights()

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the expected user types were returned
	found1 := false
	found2 := false
	for _, boardHeight := range boardHeights {
		if boardHeight.Height == boardHeight1.Height {
			found1 = true
		}
		if boardHeight.Height == boardHeight2.Height {
			found2 = true
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)

	// Cleanup
	db.Database.Unscoped().Delete(&boardHeight1)
	db.Database.Unscoped().Delete(&boardHeight2)
}
