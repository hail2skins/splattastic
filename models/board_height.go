package models

import (
	"errors"
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

// GetBoardHeights is a function that returns all board heights
func GetBoardHeights() ([]BoardHeight, error) {
	var boardHeights []BoardHeight
	result := db.Database.Find(&boardHeights)
	if result.Error != nil {
		log.Printf("Error getting board heights: %v", result.Error)
		return nil, result.Error
	}

	return boardHeights, nil
}

// BoardHeightShow is a function that returns a single board height
func BoardHeightShow(id uint64) (*BoardHeight, error) {
	var boardHeight BoardHeight
	result := db.Database.First(&boardHeight, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Board height not found")
		}
		log.Printf("Error getting board height: %v", result.Error)
		return nil, errors.New("Error getting board height")
	}

	return &boardHeight, nil
}

// Update is a method that updates a board height
func (boardheight *BoardHeight) Update(height float32) error {
	boardheight.Height = height
	result := db.Database.Save(boardheight)
	if result.Error != nil {
		log.Printf("Error updating board height: %v", result.Error)
		return result.Error
	}
	return nil
}
