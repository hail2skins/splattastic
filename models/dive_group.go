package models

import "gorm.io/gorm"

// DiveGroup struct represents the group of dives a diver can perform (Forward/Back/Referse/Inward/Twister
type DiveGroup struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}
