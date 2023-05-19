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

// StatesGet gets all states
func StatesGet() ([]State, error) {
	var states []State
	result := db.Database.Order("code ASC").Find(&states)
	if result.Error != nil {
		log.Printf("Error getting states: %v", result.Error)
		return nil, result.Error
	}
	return states, nil
}

// StateShow gets a single state
func StateShow(id uint64) (*State, error) {
	var state State
	result := db.Database.First(&state, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("state not found")
		}
		log.Printf("Error getting state: %v", result.Error)
		return nil, errors.New("error getting state")
	}
	return &state, nil
}

// Update method updates a state
func (state *State) Update(name string, code string) error {
	if name == "" || code == "" {
		return errors.New("state name and code cannot be empty")
	}

	state.Name = name
	state.Code = code

	result := db.Database.Save(state)
	if result.Error != nil {
		log.Printf("Error updating state: %v", result.Error)
		return result.Error
	}

	log.Printf("State updated: %v", state)
	return nil
}

// StateDelete deletes a state
func StateDelete(id uint64) error {
	var state State
	result := db.Database.Delete(&state, id)
	if result.Error != nil {
		log.Printf("Error deleting state: %v", result.Error)
		return result.Error
	}
	return nil
}
