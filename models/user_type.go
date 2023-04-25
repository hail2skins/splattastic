package models

import (
	"errors"
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

// UserTypeShow returns a single user type when selected
func UserTypeShow(id uint64) (*UserType, error) {
	var userType UserType
	result := db.Database.Where("id = ?", id).First(&userType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User Type not found")
		}
		log.Printf("Error getting user type: %v", result.Error)
		return nil, errors.New("Error finding User Type")
	}
	return &userType, nil
}

// Update function to update a user type
func (usertype *UserType) Update(name string) error {
	usertype.Name = name
	result := db.Database.Save(usertype)
	if result.Error != nil {
		log.Printf("Error updating User Type: %v\n", result.Error)
		return result.Error
	}
	return nil
}

// UserTypeDelete soft deletes a User Type
func UserTypeDelete(id uint64) error {
	// This only marks the deleted_at column as the current time
	result := db.Database.Where("id = ?", id).Delete(&UserType{})
	if result.Error != nil {
		log.Printf("Error deleting User Type: %v\n", result.Error)
		return result.Error
	}
	return nil
}
