package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `gorm:"not null" json:"name"`
	Location    string    `json:"location"`
	Date        time.Time `gorm:"not null" json:"date"`
	Against     string    `json:"against"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	EventTypeID uint      `gorm:"not null" json:"event_type_id"`                  // foreign key for the event type
	EventType   EventType `gorm:"constraint:OnUpdate:CASCADE;" json:"event_type"` // foreign key relationship with the EventType model without affecting deletion
}
