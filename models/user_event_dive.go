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

// DeleteUserEventDivesByEventID deletes all the dives associated with an event
func DeleteUserEventDivesByEventID(eventID uint) error {
	err := db.Database.Where("event_id = ?", eventID).Unscoped().Delete(&UserEventDive{}).Error
	if err != nil {
		return err
	}
	return nil
}
