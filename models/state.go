package models

import "gorm.io/gorm"

// State struct is the table for states
type State struct {
	gorm.Model
	Name         string `gorm:"unique;not null" json:"name"`
	Abbreviation string `gorm:"unique; not null" json:"abbreviation"`
	Teams        []Team `json:"teams"`
	Users        []User `json:"users"`
}
