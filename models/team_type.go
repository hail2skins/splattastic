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
