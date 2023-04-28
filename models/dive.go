package models

import (
	"errors"
	"log"

	db "github.com/hail2skins/splattastic/database"
	"gorm.io/gorm"
)

// Dive struct stiches together the dive group, dive type, board type and board height.
// A dive is essentially the following combination of the above:
// board.Height board.Type dive.Number+dive_type.Name(letter) dive_group.Name dive.Name dive_type.Position dive.Difficulty
// This will be a complex new page that will take lots of time and the index will be intriguing as well.
type Dive struct {
	gorm.Model
	Name          string      `gorm:"not null" json:"name"`
	Number        int         `gorm:"unique;not null" json:"number"`
	Difficulty    float32     `gorm:"not null" json:"difficulty"`
	DiveTypeID    uint64      `gorm:"not null" json:"divetype_id"`
	DiveType      DiveType    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"dive_type"`
	DiveGroupID   uint64      `gorm:"not null" json:"divegroup_id"`
	DiveGroup     DiveGroup   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"dive_group"`
	BoardTypeID   uint64      `gorm:"not null" json:"boardtype_id"`
	BoardType     BoardType   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"board_type"`
	BoardHeightID uint64      `gorm:"not null" json:"boardheight_id"`
	BoardHeight   BoardHeight `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"board_height"`
}

// DiveCreate creates a new dive
// Need to have DiveType, DiveGroup, BoardType and BoardHeight
func DiveCreate(name string, number int, difficulty float32, divetypeID uint64, divegroupID uint64, boardtypeID uint64, boardheightID uint64) (*Dive, error) {
	// Check if the associated records exist
	_, err := DiveTypeShow(divetypeID)
	if err != nil {
		log.Printf("Error getting divetype by id: %v", err)
		return nil, errors.New("Error getting divetype by id")
	}
	_, err = DiveGroupShow(divegroupID)
	if err != nil {
		log.Printf("Error getting divegroup by id: %v", err)
		return nil, errors.New("Error getting divegroup by id")
	}
	_, err = BoardTypeShow(boardtypeID)
	if err != nil {
		log.Printf("Error getting boardtype by id: %v", err)
		return nil, errors.New("Error getting boardtype by id")
	}
	_, err = BoardHeightShow(boardheightID)
	if err != nil {
		log.Printf("Error getting boardheight by id: %v", err)
		return nil, errors.New("Error getting boardheight by id")
	}

	dive := &Dive{
		Name:          name,
		Number:        number,
		Difficulty:    difficulty,
		DiveTypeID:    divetypeID,
		DiveGroupID:   divegroupID,
		BoardTypeID:   boardtypeID,
		BoardHeightID: boardheightID,
	}
	result := db.Database.Create(dive)
	if result.Error != nil {
		log.Printf("Error creating dive: %v", result.Error)
		return nil, result.Error
	}
	return dive, nil
}

// DivesGet gets all dives
func DivesGet() ([]*Dive, error) {
	dives := []*Dive{}
	// Gorm preload of associated fields
	result := db.Database.Preload("DiveType").Preload("DiveGroup").Preload("BoardType").Preload("BoardHeight").Find(&dives)
	if result.Error != nil {
		log.Printf("Error getting dives: %v", result.Error)
		return nil, result.Error
	}
	return dives, nil
}

// DiveShow gets a dive by id
func DiveShow(id uint64) (*Dive, error) {
	dive := &Dive{}
	result := db.Database.Preload("DiveType").Preload("DiveGroup").Preload("BoardType").Preload("BoardHeight").First(&dive, id)
	if result.Error != nil {
		log.Printf("Error getting dive: %v", result.Error)
		return nil, result.Error
	}
	return dive, nil
}

// DiveUpdate updates a dive
func (dive *Dive) Update(name string, number int, difficulty float32, divetypeID uint64, divegroupID uint64, boardtypeID uint64, boardheightID uint64) error {
	// Check if the associated records exist
	_, err := DiveTypeShow(divetypeID)
	if err != nil {
		log.Printf("Error getting divetype by id: %v", err)
		return errors.New("Error getting divetype by id")
	}
	_, err = DiveGroupShow(divegroupID)
	if err != nil {
		log.Printf("Error getting divegroup by id: %v", err)
		return errors.New("Error getting divegroup by id")
	}

	_, err = BoardTypeShow(boardtypeID)
	if err != nil {
		log.Printf("Error getting boardtype by id: %v", err)
		return errors.New("Error getting boardtype by id")
	}
	_, err = BoardHeightShow(boardheightID)
	if err != nil {
		log.Printf("Error getting boardheight by id: %v", err)
		return errors.New("Error getting boardheight by id")
	}

	// Update dive fields
	dive.Name = name
	dive.Number = number
	dive.Difficulty = difficulty
	dive.DiveTypeID = divetypeID
	dive.DiveGroupID = divegroupID
	dive.BoardTypeID = boardtypeID
	dive.BoardHeightID = boardheightID

	// Save updated dive to the database
	err = db.Database.Model(&dive).Updates(Dive{
		Name:          name,
		Number:        number,
		Difficulty:    difficulty,
		DiveTypeID:    divetypeID,
		DiveGroupID:   divegroupID,
		BoardTypeID:   boardtypeID,
		BoardHeightID: boardheightID,
	}).Error
	if err != nil {
		log.Printf("Error updating dive: %v", err)
		return errors.New("Error updating dive")
	}

	return nil
}
