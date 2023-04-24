package models

import (
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

type UserType struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
}

// GetUserTypes gets all user types. used to populate signup dropdown. Also for User Types Index page in admin space.
func GetUserTypes() ([]UserType, error) {
	var userTypes []UserType
	result := db.Database.Find(&userTypes)
	if result.Error != nil {
		log.Printf("Error getting user types: %v", result.Error)
		return nil, result.Error
	}
	return userTypes, nil
}

// CreateUserType creates a new user type
func CreateUserType(name string) (*UserType, error) {
	userType := UserType{Name: name}
	result := db.Database.Create(&userType)
	if result.Error != nil {
		log.Printf("Error creating user type: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Created user type: %v", userType)
	return &userType, nil
}
