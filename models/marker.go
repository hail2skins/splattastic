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

// MarkersGet is a function which will get all Markers
func MarkersGet() ([]Marker, error) {
	var markers []Marker
	result := db.Database.Find(&markers)
	if result.Error != nil {
		log.Printf("Error getting markers: %v", result.Error)
		return nil, result.Error
	}
	return markers, nil
}

// MarkerShow is a function which will get a single Marker
func MarkerShow(id uint64) (*Marker, error) {
	var marker Marker
	result := db.Database.First(&marker, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("marker not found")
		}
		log.Printf("Error getting marker: %v", result.Error)
		return nil, errors.New("error getting marker")
	}
	return &marker, nil
}

// Update method updates a marker
func (marker *Marker) Update(name string) error {
	if name == "" {
		return errors.New("marker name cannot be empty")
	}

	marker.Name = name
	result := db.Database.Save(&marker)
	if result.Error != nil {
		log.Printf("Error updating marker: %v", result.Error)
		return result.Error
	}
	return nil
}
