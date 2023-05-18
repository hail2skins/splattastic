package models

import (
	"errors"
	"log"
	"time"

	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/helpers"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string   `gorm:"unique;not null" json:"email"`
	Password   string   `gorm:"not null" json:"-"`
	UserName   string   `gorm:"unique;not null" json:"username"`
	FirstName  string   `gorm:"not null" json:"firstname"`
	LastName   string   `gorm:"not null" json:"lastname"`
	Admin      *bool    `gorm:"default:false" json:"admin"`
	UserTypeID uint64   `gorm:"not null" json:"usertype_id"`
	UserType   UserType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_type"`
	Markers    []Marker `gorm:"many2many:user_markers;association_jointable_foreignkey:marker_id;jointable_foreignkey:user_id;" json:"markers"`
	Teams      []Team   `gorm:"many2many:user_teams;association_jointable_foreignkey:team_id;jointable_foreignkey:user_id;" json:"teams"`
}

// CheckEmailUsernameAvailable checks if the email is available
func CheckEmailUsernameAvailable(email string, username string) (bool, error) {
	var user User
	result := db.Database.Where("email = ? OR user_name = ?", email, username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return true, nil
		}
		log.Printf("Error checking if email or username is available: %v", result.Error)
		return false, result.Error
	}
	return false, nil
}

// GetUserByEmail gets a user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	result := db.Database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		log.Printf("Error getting user by email: %v", result.Error)
		return nil, errors.New("Error getting user by email")
	}
	return &user, nil
}

// UserCreate creates a new user
func UserCreate(email string, password string, firstname string, lastname string, username string, usertypeName string, markerNames ...string) (*User, error) {
	hshPasswd, err := helpers.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("Error hashing password")
	}

	var usertype UserType
	if err := db.Database.Where("name = ?", usertypeName).First(&usertype).Error; err != nil {
		log.Printf("Error finding usertype: %v", err)
		return nil, errors.New("Error finding usertype")
	}

	entry := User{
		Email:     email,
		Password:  hshPasswd,
		FirstName: firstname,
		LastName:  lastname,
		UserName:  username,
		UserType:  usertype,
	}

	// Always associate the Test Marker used for testing
	markerNames = append(markerNames, "Test Marker")
	// Current marker is for Alpha ending July 1, 2023
	entry.Markers = findMarkers(markerNames, time.Now())

	result := db.Database.Create(&entry)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		return nil, errors.New("Error creating user")
	}

	// Marker for any user created before July 1, 2023 has a marker

	return &entry, nil
}

// UserFindByEmailAndPassword finds a user by email and password for the login function
func UserFindByEmailAndPassword(email string, password string) (*User, error) {
	var user User
	result := db.Database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		log.Printf("Error getting user by email: %v", result.Error)
		return nil, errors.New("Error getting user by email")
	}

	match := helpers.CheckPasswordHash(password, user.Password)
	if match {
		//log.Printf("UserFindByEmailAndPassword: user = %v", user) // for debugging
		return &user, nil
	} else {
		return nil, errors.New("Password does not match")
	}
}

// UserFind finds a user by id
func UserFind(id uint64) (*User, error) {
	var user User
	result := db.Database.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		log.Printf("Error getting user by id: %v", result.Error)
		return nil, errors.New("Error finding user")
	}
	return &user, nil
}

// UserShow shows a user by id and associated UserType for a profile page
func UserShow(id uint64) (*User, error) {
	var user User
	result := db.Database.Preload("UserType").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		log.Printf("Error getting user by id: %v", result.Error)
		return nil, errors.New("Error finding user")
	}
	return &user, nil
}

// Update method
func (user *User) Update(email string, firstname string, lastname string, username string, usertypeID uint64) error {
	_, err := UserTypeShow(usertypeID)
	if err != nil {
		log.Printf("Error finding usertype: %v", err)
		return errors.New("Error finding usertype")
	}

	// Update user fields
	user.Email = email
	user.FirstName = firstname
	user.LastName = lastname
	user.UserName = username
	user.UserTypeID = usertypeID

	// Save updated user to the db
	err = db.Database.Model(&user).Updates(User{
		Email:      email,
		FirstName:  firstname,
		LastName:   lastname,
		UserName:   username,
		UserTypeID: usertypeID,
	}).Error
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return errors.New("Error updating user")
	}
	return nil

}
