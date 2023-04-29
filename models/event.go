package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name     string    `gorm:"not null" json:"name"`                                      // really more a description than name
	Location string    `gorm:"not null" json:"location"`                                  // eventually will use addresses here but now just name of school or whatever
	Date     time.Time `gorm:"not null" json:"date"`                                      // date selected through form datepicker can be forward or backward in time
	Against  string    `gorm:"not null" json:"against"`                                   // name of opposing team or school
	UserID   uint      `gorm:"not null" json:"user_id"`                                   // foreign key for the user who created the event
	User     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"` // foreign key relationship with the User model
}
