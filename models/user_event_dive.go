package models

import "gorm.io/gorm"

// UserEventDive struct is a join table between UserEvent and Dive
type UserEventDive struct {
	gorm.Model
	UserID  uint `gorm:"primaryKey"`
	EventID uint `gorm:"primaryKey"`
	DiveID  uint `gorm:"primaryKey"`

	User  User
	Event Event
	Dive  Dive
}
