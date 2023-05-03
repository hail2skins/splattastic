package helpers

import (
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func DeleteUserEventDives() error {
	var userEventDives []models.UserEventDive
	result := db.Database.Find(&userEventDives)
	if result.Error != nil {
		return result.Error
	}

	//fmt.Println("UserEventDives before deletion:", userEventDives)

	for _, userEventDive := range userEventDives {
		db.Database.Unscoped().Delete(&userEventDive)
	}

	// Check if the user event dives are deleted
	var userEventDivesAfterDeletion []models.UserEventDive
	result = db.Database.Find(&userEventDivesAfterDeletion)
	if result.Error != nil {
		return result.Error
	}
	//fmt.Println("UserEventDives after deletion:", userEventDivesAfterDeletion)

	return nil
}
