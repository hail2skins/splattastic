package models

import "gorm.io/gorm"

type UserType string

const (
	Athlete   UserType = "Athlete"
	Coach     UserType = "Coach"
	Owner     UserType = "Owner"
	Supporter UserType = "Supporter"
)

type User struct {
	gorm.Model
	Email     string   `gorm:"unique;not null" json:"email"`
	Password  string   `gorm:"not null" json:"-"`
	UserName  string   `gorm:"unique;not null" json:"username"`
	FirstName string   `gorm:"not null" json:"firstname"`
	LastName  string   `gorm:"not null" json:"lastname"`
	Admin     bool     `gorm:"default:false" json:"admin"`
	UserType  UserType `gorm:"type:ENUM('Athlete', 'Coach', 'Owner', 'Supporter');not null" json:"usertype"`
}
