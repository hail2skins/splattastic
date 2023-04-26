package models

import "gorm.io/gorm"

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
