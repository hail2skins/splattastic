package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// EventType is a model which will eventually tie to the type of event.
// An EventType is beleived to only be either a "practice" or a "meet" but we could change and add something later.
type EventType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}

// EventTypeCreate is a function which will create a new EventType
func EventTypeCreate(name string) (*EventType, error) {
	if name == "" {
		return nil, errors.New("event type name cannot be empty")
	}

	eventType := EventType{Name: name}
	result := db.Database.Create(&eventType)
	if result.Error != nil {
		log.Printf("Error creating event type: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("Event type created: %v", eventType)
	return &eventType, nil
}

// EventTypesGet is a function which will get all EventTypes
func EventTypesGet() ([]EventType, error) {
	var eventTypes []EventType
	result := db.Database.Find(&eventTypes)
	if result.Error != nil {
		log.Printf("Error getting event types: %v", result.Error)
		return nil, result.Error
	}
	return eventTypes, nil
}

// EventTypeShow is a function which will get a single EventType
func EventTypeShow(id uint64) (*EventType, error) {
	var eventType EventType
	result := db.Database.First(&eventType, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("event type not found")
		}
		log.Printf("Error getting event type: %v", result.Error)
		return nil, errors.New("error getting event type")
	}
	return &eventType, nil
}

// Update method updates a event type
func (eventType *EventType) Update(name string) error {
	eventType.Name = name
	result := db.Database.Save(&eventType)
	if result.Error != nil {
		log.Printf("Error updating event type: %v", result.Error)
		return result.Error
	}
	return nil
}

// EventTypeDelete is a function which will delete a single EventType
func EventTypeDelete(id uint64) error {
	var eventType EventType
	result := db.Database.Delete(&eventType, id)
	if result.Error != nil {
		log.Printf("Error deleting event type: %v", result.Error)
		return result.Error
	}
	return nil
}
