package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// TeamType is the model for the team_type table representing the type of team
// which is a model to come. Right now the teams should be high school or club
// but we'll add college and others as we go.
type TeamType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}

// TeamTypeCreate is a function which will create a new TeamType
func TeamTypeCreate(name string) (*TeamType, error) {
	if name == "" {
		return nil, errors.New("team type name cannot be empty")
	}

	teamType := TeamType{Name: name}
	result := db.Database.Create(&teamType)
	if result.Error != nil {
		log.Printf("Error creating team type: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("Team type created: %v", teamType)
	return &teamType, nil
}

// TeamTypesGet is a function which will get all TeamTypes
func TeamTypesGet() ([]TeamType, error) {
	var teamTypes []TeamType
	result := db.Database.Find(&teamTypes)
	if result.Error != nil {
		log.Printf("Error getting team types: %v", result.Error)
		return nil, result.Error
	}
	return teamTypes, nil
}

// TeamTypeShow is a function which will get a single TeamType
func TeamTypeShow(id uint64) (*TeamType, error) {
	var teamType TeamType
	result := db.Database.First(&teamType, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("team type not found")
		}
		log.Printf("Error getting team type: %v", result.Error)
		return nil, errors.New("error getting team type")
	}
	return &teamType, nil
}
