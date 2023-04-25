package models

import (
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// BoardHeight struct represents the height of the boards/platforms a diver can jump from
type BoardHeight struct {
	gorm.Model
	Height float32 `gorm:"unique;not null" json:"height"`
}

// CreateBoardHeight is a function that creates a new board height
func CreateBoardHeight(height float32) (*BoardHeight, error) {
	boardHeight := BoardHeight{Height: height}
	result := db.Database.Create(&boardHeight)
	if result.Error != nil {
		log.Printf("Error creating board height: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Board height created: %v", boardHeight)
	return &boardHeight, nil
}
