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
	Name   string `gorm:"not null;unique" json:"name"`   // Position of the dive (Straight/Pike/Tuck/Free)
	Letter string `gorm:"not null;unique" json:"letter"` // Letter of the dive (A/B/C/D)
}

// DiveTypeCreate creates a new dive type
func DiveTypeCreate(name string, letter string) (*DiveType, error) {
	if name == "" {
		return nil, errors.New("dive type name cannot be empty")
	}

	diveType := DiveType{Name: name, Letter: letter}
	result := db.Database.Create(&diveType)
	if result.Error != nil {
		log.Printf("Error creating dive type: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Dive type created: %v", diveType)
	return &diveType, nil
}

// DiveTypesGet gets all dive types
func DiveTypesGet() ([]DiveType, error) {
	var diveTypes []DiveType
	result := db.Database.Order("letter ASC").Find(&diveTypes)
	if result.Error != nil {
		log.Printf("Error getting dive types: %v", result.Error)
		return nil, result.Error
	}
	return diveTypes, nil
}

// DiveTypeShow gets a single dive type
func DiveTypeShow(id uint64) (*DiveType, error) {
	var diveType DiveType
	result := db.Database.First(&diveType, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("dive type not found")
		}
		log.Printf("Error getting dive type: %v", result.Error)
		return nil, errors.New("error getting dive type")
	}
	return &diveType, nil
}

// Update method updates a dive type
func (diveType *DiveType) Update(name string, letter string) error {
	diveType.Name = name
	diveType.Letter = letter
	result := db.Database.Save(diveType)
	if result.Error != nil {
		log.Printf("Error updating dive type: %v", result.Error)
		return result.Error
	}
	return nil
}

// DiveTypeDelete deletes a dive type
func DiveTypeDelete(id uint64) error {
	var diveType DiveType
	result := db.Database.Delete(&diveType, id)
	if result.Error != nil {
		log.Printf("Error deleting dive type: %v", result.Error)
		return result.Error
	}
	return nil
}
