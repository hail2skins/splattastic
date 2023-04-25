package models

import "gorm.io/gorm"

// BoardType is a model for the board_types table
// A board type in diving is just springboard or platform.
// In theory we could add "cliff" if that seemed needed :).
type BoardType struct {
	gorm.Model
	Name          string      `gorm:"unique;not null" json:"name" form:"name" binding:"required"`
	BoardHeightID uint64      `gorm:"not null" json:"board_height_id"`
	BoardHeight   BoardHeight `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_type"`
}
