package models

import (
	"errors"
	"log"
	"regexp"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// Team struct is is the table for teams
type Team struct {
	gorm.Model
	Name         string   `gorm:"unique;not null" json:"name"`
	Street       string   `json:"street"`
	Street1      string   `json:"street1"`
	City         string   `json:"city"`
	Zip          string   `gorm:"not null" json:"zip"`
	Abbreviation string   `json:"abbreviation"`
	TeamTypeID   uint64   `gorm:"not null" json:"teamtype_id"`
	TeamType     TeamType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"team_type"`
	Users        []User   `gorm:"many2many:user_teams;association_jointable_foreignkey:user_id;jointable_foreignkey:team_id;" json:"users"`
	StateID      uint64   `gorm:"not null" json:"state_id"`
	State        State    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"state"`
}

// TeamCreate creates a new team
// Needs the FK relationship with team_type and needs to
// be associated with a state. Also needs lots of valication
func TeamCreate(name string, street string, street1 string, city string, zip string, abbreviation string, teamTypeID uint64, stateID uint64) (*Team, error) {
	// Check if the associated records exist
	_, err := TeamTypeShow(teamTypeID)
	if err != nil {
		return nil, err
	}
	_, err = StateShow(stateID)
	if err != nil {
		return nil, err
	}

	// validate zip code
	zipRegex := `^\d{5}$|^\d{5}-\d{4}$`
	match, _ := regexp.MatchString(zipRegex, zip)
	if !match {
		return nil, errors.New("Invalid zip code must be in the format 12345 or 12345-1234")
	}

	team := &Team{
		Name:         name,
		Street:       street,
		Street1:      street1,
		City:         city,
		Zip:          zip,
		Abbreviation: abbreviation,
		TeamTypeID:   teamTypeID,
		StateID:      stateID,
	}

	result := db.Database.Create(team)
	if result.Error != nil {
		log.Printf("Error creating team: %v", result.Error)
		return nil, result.Error
	}

	return team, nil
}
