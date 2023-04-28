package helpers

import (
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// CleanTestData deletes the test data created by CreateTestData.
func CleanTestData(dg1, dg2 *models.DiveGroup, dt1, dt2 *models.DiveType, bt1, bt2 *models.BoardType, bh1, bh2 *models.BoardHeight) {
	db.Database.Unscoped().Delete(&models.DiveGroup{}, dg1.ID)
	db.Database.Unscoped().Delete(&models.DiveGroup{}, dg2.ID)
	db.Database.Unscoped().Delete(&models.DiveType{}, dt1.ID)
	db.Database.Unscoped().Delete(&models.DiveType{}, dt2.ID)
	db.Database.Unscoped().Delete(&models.BoardType{}, bt1.ID)
	db.Database.Unscoped().Delete(&models.BoardType{}, bt2.ID)
	db.Database.Unscoped().Delete(&models.BoardHeight{}, bh1.ID)
	db.Database.Unscoped().Delete(&models.BoardHeight{}, bh2.ID)
}
