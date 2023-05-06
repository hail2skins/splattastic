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
	//log.Printf("Dive IDs in EventCreate: %v", diveIDs) // Add this line for debugging
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
		result := db.Database.Create(userEventDive)
		if result.Error != nil {
			log.Printf("Error creating UserEventDive: %v", result.Error)
		} else {
			log.Printf("Created UserEventDive with ID: %d", userEventDive.ID)
		}
	}

	return event, nil
}

// HasDive method is really for the model test.  May want to move this a bit later, though it may come in handy here we'll see.
func (e *Event) HasDive(dive *Dive) bool {
	var userEventDives []UserEventDive
	db.Database.Where("event_id = ?", e.ID).Find(&userEventDives)

	//fmt.Printf("UserEventDives: %+v\n", userEventDives) // Add this line for debugging

	// loop through the userEventDives and see if the diveID matches the dive.ID
	for _, userEventDive := range userEventDives {
		if userEventDive.DiveID == uint64(dive.ID) {
			return true
		}
	}

	return false
}

// EventShow returns a single event with related dives
func EventShow(id uint64) (*Event, error) {
	event := &Event{}
	result := db.Database.Preload("User").Preload("EventType").First(&event, id)
	if result.Error != nil {
		log.Printf("Error retrieving event: %v", result.Error)
		return nil, result.Error
	}

	// Find the associated UserEventDives for the event
	var userEventDives []UserEventDive
	result = db.Database.Where("event_id = ?", id).Find(&userEventDives)
	if result.Error != nil {
		log.Printf("Error retrieving userEventDives: %v", result.Error)
		return nil, result.Error
	}

	// Find the associated dives for each UserEventDive
	for _, userEventDive := range userEventDives {
		dive := &Dive{}
		result := db.Database.First(&dive, userEventDive.DiveID)
		if result.Error != nil {
			log.Printf("Error retrieving dive: %v", result.Error)
			return nil, result.Error
		}

	}

	return event, nil
}

// GetDivesForEvent retrieves the dives associated with an event
func GetDivesForEvent(eventID uint64) ([]Dive, error) {
	var userEventDives []UserEventDive
	var dives []Dive

	err := db.Database.Where("event_id = ? AND deleted_at is NULL", eventID).Find(&userEventDives).Error
	if err != nil {
		log.Printf("Error retrieving UserEventDives: %v", err)
		return nil, err
	}

	for _, userEventDive := range userEventDives {
		var dive Dive
		err = db.Database.Preload("DiveType").Preload("DiveGroup").Preload("BoardHeight").Preload("BoardType").First(&dive, userEventDive.DiveID).Error
		if err != nil {
			log.Printf("Error retrieving dive: %v", err)
			return nil, err
		}
		dives = append(dives, dive)
	}

	return dives, nil
}

// GetUserEvents retrieves all events for a user.
// We are NOT populating UserEventDives here.  Instead when needed
// in the controller we will call GetDivesForEvent to populate the
// dives for each event.
func GetUserEvents(userID uint64) ([]Event, error) {
	var events []Event
	err := db.Database.Preload("EventType").Where("user_id = ?", userID).Find(&events).Error
	if err != nil {
		log.Printf("Error retrieving events: %v", err)
		return nil, err
	}

	return events, nil
}

// GetLastFiveEvents retrieves the last five events for a user.
// We are NOT populating UserEventDives here.  Instead when needed
// in the controller we will call GetDivesForEvent to populate the
// dives for each event.
func GetLastFiveEvents(userID uint64) (*[]Event, error) {
	var events []Event
	result := db.Database.Preload("EventType").Where("user_id = ?", userID).Order("updated_at desc").Limit(5).Find(&events)
	if result.Error != nil {
		log.Printf("Error retrieving events: %v", result.Error)
		return nil, result.Error
	}

	return &events, nil
}

// Update method updates an event
func (e *Event) Update(name string, location string, date time.Time, against string, eventTypeID uint64, diveIDs []uint64) (*Event, error) {
	// Check if the associated records exist
	_, err := EventTypeShow(eventTypeID)
	if err != nil {
		return nil, err
	}

	e.Name = name
	e.Location = location
	e.Date = date
	e.Against = against
	e.EventTypeID = eventTypeID

	result := db.Database.Save(e)
	if result.Error != nil {
		log.Printf("Error updating event: %v", result.Error)
		return nil, result.Error
	}

	// Retrieve the existing dives associated with the event
	existingDives, err := GetDivesForEvent(uint64(e.ID))
	if err != nil {
		return nil, err
	}

	// Create a map of the existing dive IDs
	existingDiveIDs := make(map[uint64]bool)
	for _, dive := range existingDives {
		existingDiveIDs[uint64(dive.ID)] = true
	}

	// Add and update the dives associated with the event
	for _, diveID := range diveIDs {
		if existingDiveIDs[diveID] {
			// Dive is already associated with the event, remove it from the existingDiveIDs map
			delete(existingDiveIDs, diveID)
		} else {
			// Dive is not associated with the event, create a new UserEventDive record
			userEventDive := &UserEventDive{
				UserID:  e.UserID,
				EventID: uint64(e.ID),
				DiveID:  diveID,
			}
			result := db.Database.Create(userEventDive)
			if result.Error != nil {
				log.Printf("Error creating UserEventDive: %v", result.Error)
			} else {
				log.Printf("Created UserEventDive with ID: %d", userEventDive.ID)
			}
		}
	}

	// Remove the dives that are no longer associated with the event
	for diveID := range existingDiveIDs {
		result := db.Database.Delete(&UserEventDive{}, "user_id = ? AND event_id = ? AND dive_id = ?", e.UserID, uint64(e.ID), diveID)
		if result.Error != nil {
			log.Printf("Error deleting UserEventDive: %v", result.Error)
		} else {
			log.Printf("Deleted UserEventDive with UserID: %d, EventID: %d, and DiveID: %d", e.UserID, uint64(e.ID), diveID)
		}
	}

	return e, nil
}

// EventDelete soft deletes an event
func EventDelete(id uint64) error {
	UserEventDiveDelete(id)
	result := db.Database.Delete(&Event{}, id)
	if result.Error != nil {
		log.Printf("Error deleting event: %v", result.Error)
		return result.Error
	}

	return nil
}
