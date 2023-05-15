package models

import "gorm.io/gorm"

// UserMarker struct to store user marker data
type Marker struct {
	gorm.Model
	Name string `json:"name"`
}
