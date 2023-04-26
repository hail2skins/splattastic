package models

import "gorm.io/gorm"

// DiveType is a model for the dive_types table
type DiveType struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"` // Position of the dive (Straight/Pike/Tuck/Free)
}
