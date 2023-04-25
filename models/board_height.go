package models

import "gorm.io/gorm"

// BoardHeight struct represents the height of the boards/platforms a diver can jump from
type BoardHeight struct {
	gorm.Model
	Height float32
}
