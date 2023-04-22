package models

import "gorm.io/gorm"

type UserType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}
