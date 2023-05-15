package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// UserMarker struct to store user marker data
type Marker struct {
	gorm.Model
	Name string `gorm:"unique;notnull" json:"name"`
}

// MarkerCreate is a function which will create a new Marker
func MarkerCreate(name string) (*Marker, error) {
	if name == "" {
		return nil, errors.New("marker name cannot be empty")
	}

	marker := Marker{Name: name}
	result := db.Database.Create(&marker)
	if result.Error != nil {
		log.Printf("Error creating marker: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("Marker created: %v", marker)
	return &marker, nil
}
