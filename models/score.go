package models

import (
	"errors"
	"fmt"
	"math"
	"sort"

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

// ScoreUpsert creates a new score record or updates an existing one
func ScoreUpsert(userID uint64, eventID uint64, diveID uint64, judge int, value float64) (*Score, error) {
	// Validation.  The math.Mod check ensures that the value is in increments of 0.5
	if value < 0 || value > 10 || math.Mod(value*2, 1) != 0 {
		return nil, errors.New("Invalid score value. Score must be between 0 and 10 and in increments of 0.5.")
	}

	// Check if a score from the same judge for the same dive already exists
	var existingScore Score
	err := db.Database.Where("user_id = ? AND event_id = ? AND dive_id = ? AND judge = ?", userID, eventID, diveID, judge).First(&existingScore).Error

	if err == gorm.ErrRecordNotFound {
		// Create a new score
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
	} else if err != nil {
		// If the error is not gorm.ErrRecordNotFound, then there was a different error
		return nil, err
	} else {
		// Update the existing score
		existingScore.Value = value
		err = db.Database.Save(&existingScore).Error
		if err != nil {
			return nil, err
		}
		return &existingScore, nil
	}
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

// CalculateDiveScore calculates the total score for a dive based on FINA dive rules
// We may want to move total score to a separate table for retrieval and use elsewhere in the application
// But this is called from the EventShow page for now as dives are entered using the JS.
var ErrInvalidScores = errors.New("Invalid number of scores for dive")

func CalculateDiveScore(diveID uint64) (float64, error) {
	var scores []Score
	var dive Dive

	if err := db.Database.Where("dive_id = ?", diveID).Find(&scores).Error; err != nil {
		return 0, err
	}

	if err := db.Database.First(&dive, diveID).Error; err != nil {
		return 0, err
	}

	// Sort the scores
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Value < scores[j].Value
	})

	var totalScore float64
	switch len(scores) {
	case 3:
		for _, score := range scores {
			totalScore += score.Value
		}
	case 5:
		for _, score := range scores[1:4] {
			totalScore += score.Value
		}
	case 7:
		for _, score := range scores[2:5] {
			totalScore += score.Value
		}
	default:
		return 0, ErrInvalidScores
	}

	return totalScore * float64(dive.Difficulty), nil
}

func CalculateMeetScore(eventID uint64) (float64, error) {
	// Get all dives for this user and event
	dives, err := GetDivesForEvent(eventID)
	if err != nil {
		return 0, err
	}

	var totalScore float64 = 0
	for _, dive := range dives {
		diveScore, err := CalculateDiveScore(uint64(dive.ID))
		if err != nil {
			if err == ErrInvalidScores {
				continue // Skip this dive
			}
			return 0, err
		}
		totalScore += diveScore
	}

	return totalScore, nil
}
