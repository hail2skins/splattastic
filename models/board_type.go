package models

import "gorm.io/gorm"

// BoardType is a model for the board_types table
// A board type in diving is just springboard or platform.
// In theory we could add "cliff" if that seemed needed :).
// BoardType is a many2many relationship with BoardHeight
type BoardType struct {
	gorm.Model
	Name         string        `gorm:"unique;not null" json:"name" form:"name" binding:"required"`
	BoardHeights []BoardHeight `gorm:"many2many:board_type_board_heights;"`
}
