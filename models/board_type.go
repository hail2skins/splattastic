package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// BoardType is a model for the board_types table
// A board type in diving is just springboard or platform.
// In theory we could add "cliff" if that seemed needed :).
// BoardType is a many2many relationship with BoardHeight
type BoardType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name" form:"name" binding:"required"`
}

// BoardTypeCreate creates a new board type
func BoardTypeCreate(name string) (*BoardType, error) {
	if name == "" {
		return nil, errors.New("board type name cannot be empty")
	}
	boardType := BoardType{Name: name}
	result := db.Database.Create(&boardType)
	if result.Error != nil {
		log.Printf("Error creating board type: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Board type created: %v", boardType)
	return &boardType, nil
}

// BoardTypesGet gets all board types
func BoardTypesGet() ([]BoardType, error) {
	var boardTypes []BoardType
	result := db.Database.Find(&boardTypes)
	if result.Error != nil {
		log.Printf("Error getting board types: %v", result.Error)
		return nil, result.Error
	}
	return boardTypes, nil
}

// BoardTypeShow gets a single board type
func BoardTypeShow(id uint64) (*BoardType, error) {
	var boardType BoardType
	result := db.Database.First(&boardType, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("board type not found")
		}
		log.Printf("Error getting board type: %v", result.Error)
		return nil, errors.New("error getting board type")
	}
	return &boardType, nil
}

// Update method updates a board type
func (boardType *BoardType) Update(name string) error {
	boardType.Name = name
	result := db.Database.Save(boardType)
	if result.Error != nil {
		log.Printf("Error updating board type: %v", result.Error)
		return result.Error
	}
	return nil

}

// BoardTypeDelete soft deletes a board type
func BoardTypeDelete(id uint64) error {
	result := db.Database.Delete(&BoardType{}, id)
	if result.Error != nil {
		log.Printf("Error deleting board type: %v", result.Error)
		return result.Error
	}
	return nil
}
