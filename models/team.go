package models

import "gorm.io/gorm"

// Team struct is is the table for teams
type Team struct {
	gorm.Model
	Name         string   `gorm:"unique;not null" json:"name"`
	Street       string   `json:"street"`
	Street1      string   `json:"street1"`
	City         string   `json:"city"`
	Zip          string   `json:"zip"`
	Abbreviation string   `json:"abbreviation"`
	TeamTypeID   uint64   `gorm:"not null" json:"teamtype_id"`
	TeamType     TeamType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"team_type"`
	Users        []User   `gorm:"many2many:user_teams;association_jointable_foreignkey:user_id;jointable_foreignkey:team_id;" json:"users"`
	StateID      uint64   `gorm:"not null" json:"state_id"`
	State        State    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"state"`
}
