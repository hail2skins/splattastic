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

// DiveGroupsGet gets all dive groups
func DiveGroupsGet() ([]DiveGroup, error) {
	var diveGroups []DiveGroup
	result := db.Database.Find(&diveGroups)
	if result.Error != nil {
		log.Printf("Error getting dive groups: %v", result.Error)
		return nil, result.Error
	}
	return diveGroups, nil
}

// DiveGroupShow gets a single dive group
func DiveGroupShow(id uint64) (*DiveGroup, error) {
	var diveGroup DiveGroup
	result := db.Database.First(&diveGroup, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("dive group not found")
		}
		log.Printf("Error getting dive group: %v", result.Error)
		return nil, errors.New("error getting dive group")
	}
	return &diveGroup, nil
}
