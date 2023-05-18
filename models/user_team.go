package models

// UserTeam struct is the join table for Users and Teams
type UserTeam struct {
	UserID uint64 `gorm:"primary_key" json:"user_id"`
	TeamID uint64 `gorm:"primary_key" json:"team_id"`
}
