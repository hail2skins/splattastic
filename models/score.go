package models

import (
	"errors"
	"fmt"
	"math"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

type Score struct {
	gorm.Model
	UserID        uint64        `json:"user_id"`
	EventID       uint64        `json:"event_id"`
	DiveID        uint64        `json:"dive_id"`
	UserEventDive UserEventDive `gorm:"foreignKey:UserID,EventID,DiveID;references:UserID,EventID,DiveID" json:"user_event_dive"`
	Judge         int           `json:"judge"`
	Value         float64       `json:"score"`
}

// ScoreCreate creates a score record can have between 1 and 9 judges
func ScoreCreate(userID uint64, eventID uint64, diveID uint64, judge int, value float64) (*Score, error) {
	// Validation.  The math.Mod check ensures that the value is in increments of 0.5
	if value < 0 || value > 10 || math.Mod(value*2, 1) != 0 {
		return nil, errors.New("Invalid score value. Score must be between 0 and 10 and in increments of 0.5.")
	}
	score := Score{
		UserID:  userID,
		EventID: eventID,
		DiveID:  diveID,
		Judge:   judge,
		Value:   value,
	}
	err := db.Database.Create(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

// FetchScores retrieves all the scores for a specific event for use with JS
func FetchScores(userID uint64, eventID uint64, diveID uint64) ([]Score, error) {
	var scores []Score
	err := db.Database.Where("user_id = ? AND event_id = ? AND dive_id = ?", userID, eventID, diveID).Find(&scores).Error
	if err != nil {
		return nil, err
	}

	// Debugging: Print the fetched scores
	fmt.Printf("Fetched scores for user %d, event %d, dive %d: %+v\n", userID, eventID, diveID, scores)

	return scores, nil
}
