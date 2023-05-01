package models

import (
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
