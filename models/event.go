package models

import (
	"log"
	"time"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `gorm:"not null" json:"name"`
	Location    string    `json:"location"`
	Date        time.Time `gorm:"not null" json:"date"`
	Against     string    `json:"against"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	EventTypeID uint64    `gorm:"not null" json:"event_type_id"`                  // foreign key for the event type
	EventType   EventType `gorm:"constraint:OnUpdate:CASCADE;" json:"event_type"` // foreign key relationship with the EventType model without affecting deletion
}

// EventCreate creates a new event
// Need the fk relationship wit the user and event type
// Dives may have been selected.  Between and many.  Need to figure out how to do that.
func EventCreate(name string, location string, date time.Time, against string, userID uint64, eventTypeID uint64, diveIDs []uint64) (*Event, error) {
	// Check if the associated records exist
	_, err := UserShow(userID)
	if err != nil {
		return nil, err
	}
	_, err = EventTypeShow(eventTypeID)
	if err != nil {
		return nil, err
	}

	event := &Event{
		Name:        name,
		Location:    location,
		Date:        date,
		Against:     against,
		UserID:      userID,
		EventTypeID: eventTypeID,
	}

	result := db.Database.Create(event)
	if result.Error != nil {
		log.Printf("Error creating event: %v", result.Error)
		return nil, result.Error
	}

	// Add the dives to the event
	for _, diveID := range diveIDs {
		userEventDive := &UserEventDive{
			UserID:  userID,
			EventID: uint64(event.ID),
			DiveID:  diveID,
		}
		db.Database.Create(userEventDive)
	}

	return event, nil
}

func (e *Event) HasDive(dive *Dive) bool {
	var userEventDives []UserEventDive
	db.Database.Where("event_id = ?", e.ID).Find(&userEventDives)

	//fmt.Printf("UserEventDives: %+v\n", userEventDives) // Add this line for debugging

	for _, userEventDive := range userEventDives {
		if userEventDive.DiveID == uint64(dive.ID) {
			return true
		}
	}

	return false
}
