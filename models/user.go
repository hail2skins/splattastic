package models

import (
	"errors"
	"log"

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
	Admin      bool     `gorm:"default:false" json:"admin"`
	UserTypeID uint64   `gorm:"not null" json:"usertype_id"`
	UserType   UserType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_type"`
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
func UserCreate(email string, password string, username string, firstname string, lastname string, usertypeName string) (*User, error) {
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
		UserName:  username,
		FirstName: firstname,
		LastName:  lastname,
		UserType:  usertype,
	}

	result := db.Database.Create(&entry)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		return nil, errors.New("Error creating user")
	}

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
