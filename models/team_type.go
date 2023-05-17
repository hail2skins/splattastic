package models

import "gorm.io/gorm"

// TeamType is the model for the team_type table representing the type of team
// which is a model to come. Right now the teams should be high school or club
// but we'll add college and others as we go.
type TeamType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}
