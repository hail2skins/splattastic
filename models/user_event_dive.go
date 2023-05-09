package models

import (
	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// UserEventDive struct is a join table between UserEvent and Dive
type UserEventDive struct {
	gorm.Model
	UserID  uint64 `gorm:"primaryKey"`
	EventID uint64 `gorm:"primaryKey"`
	DiveID  uint64 `gorm:"primaryKey"`

	User  User
	Event Event
	Dive  Dive
}

// DeleteUserEventDivesByEventID deletes all the dives associated with an event in permanent fashion
func DeleteUserEventDivesByEventID(eventID uint) error {
	err := db.Database.Where("event_id = ?", eventID).Unscoped().Delete(&UserEventDive{}).Error
	if err != nil {
		return err
	}
	return nil
}

// UserEventDiveDelete deletes a UserEventDive record soft delete for Update and probably Delete event functions
func UserEventDiveDelete(eventID uint64) error {
	err := db.Database.Where("event_id = ?", eventID).Delete(&UserEventDive{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserEventDivesForEvent retrieves UserEventDive records for a specific event
func GetUserEventDivesForEvent(eventID uint64) ([]UserEventDive, error) {
	var userEventDives []UserEventDive
	err := db.Database.Where("event_id = ?", eventID).Find(&userEventDives).Error
	return userEventDives, err
}
