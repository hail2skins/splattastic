package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// State struct is the table for states
type State struct {
	gorm.Model
	Name  string `gorm:"unique;not null" json:"name"`
	Code  string `gorm:"unique; not null" json:"code"`
	Teams []Team `json:"teams"`
	Users []User `json:"users"`
}

// StateCreate creates a new state
func StateCreate(name, code string) (*State, error) {
	if name == "" || code == "" {
		return nil, errors.New("state name and code cannot be empty")
	}

	state := State{
		Name: name,
		Code: code,
	}

	result := db.Database.Create(&state)
	if result.Error != nil {
		log.Printf("Error creating state: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("State created: %v", state)
	return &state, nil
}
