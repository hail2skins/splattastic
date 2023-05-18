package models

import "gorm.io/gorm"

// State struct is the table for states
type State struct {
	gorm.Model
	Name  string `gorm:"unique;not null" json:"name"`
	Code  string `gorm:"unique; not null" json:"code"`
	Teams []Team `json:"teams"`
	Users []User `json:"users"`
}
