package models

// UserMarker struct is the join table for Users and Markers
type UserMarker struct {
	UserID   uint64 `gorm:"primary_key" json:"user_id"`
	MarkerID uint64 `gorm:"primary_key" json:"marker_id"`
}
