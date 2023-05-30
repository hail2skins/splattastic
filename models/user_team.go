package models

import (
	"log"

	db "github.com/hail2skins/splattastic/database"
)

// UserTeam struct is the join table for Users and Teams
type UserTeam struct {
	UserID uint64 `gorm:"primary_key" json:"user_id"`
	TeamID uint64 `gorm:"primary_key" json:"team_id"`
}

// UserTeamCreate creates a new user_team record for Users and Teams
func UserTeamCreate(userID uint64, teamID uint64) error {
	userTeam := &UserTeam{
		UserID: userID,
		TeamID: teamID,
	}
	result := db.Database.Create(userTeam)
	if result.Error != nil {
		log.Printf("Creating user_team association with User ID: %d and Team ID: %d", userID, teamID)
		log.Printf("Error creating user_team: %v", result.Error)
		return result.Error
	}
	return nil
}
