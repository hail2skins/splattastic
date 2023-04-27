package helpers

import (
	"github.com/hail2skins/splattastic/models"
)

// CreateTestData creates test data for DiveGroup, DiveType, BoardType, and BoardHeight.
func CreateTestData() (dg1, dg2 *models.DiveGroup, dt1, dt2 *models.DiveType, bt1, bt2 *models.BoardType, bh1, bh2 *models.BoardHeight) {
	dg1, _ = models.DiveGroupCreate("Test Dive Group 1")
	dg2, _ = models.DiveGroupCreate("Test Dive Group 2")
	dt1, _ = models.DiveTypeCreate("Test Dive Type 1")
	dt2, _ = models.DiveTypeCreate("Test Dive Type 2")
	bt1, _ = models.BoardTypeCreate("Test Board Type 1")
	bt2, _ = models.BoardTypeCreate("Test Board Type 2")
	bh1, _ = models.CreateBoardHeight(8.5)
	bh2, _ = models.CreateBoardHeight(6)

	return
}
