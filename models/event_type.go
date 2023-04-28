package models

import "gorm.io/gorm"

// EventType is a model which will eventually tie to the type of event.
// An EventType is beleived to only be either a "practice" or a "meet" but we could change and add something later.
type EventType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}
