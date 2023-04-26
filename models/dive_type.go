package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// DiveType is a model for the dive_types table
type DiveType struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"` // Position of the dive (Straight/Pike/Tuck/Free)
}

// DiveTypeCreate creates a new dive type
func DiveTypeCreate(name string) (*DiveType, error) {
	if name == "" {
		return nil, errors.New("dive type name cannot be empty")
	}

	diveType := DiveType{Name: name}
	result := db.Database.Create(&diveType)
	if result.Error != nil {
		log.Printf("Error creating dive type: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Dive type created: %v", diveType)
	return &diveType, nil
}
