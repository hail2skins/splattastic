package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// DiveGroup struct represents the group of dives a diver can perform (Forward/Back/Referse/Inward/Twister
type DiveGroup struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}

// DiveGroupCreate creates a new dive group
func DiveGroupCreate(name string) (*DiveGroup, error) {
	if name == "" {
		return nil, errors.New("dive group name cannot be empty")
	}

	diveGroup := DiveGroup{Name: name}
	result := db.Database.Create(&diveGroup)
	if result.Error != nil {
		log.Printf("Error creating dive group: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Dive group created: %v", diveGroup)
	return &diveGroup, nil
}
