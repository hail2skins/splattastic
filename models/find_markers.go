package models

import (
	"log"
	"time"

	db "github.com/hail2skins/splattastic/database"
)

// FindMarkers finds and creates any active markers.
// the initial code ensures we can test the marker function
// even when no markers are in effect.  This is called within
// the UserCreate function.
func findMarkers(markerNames []string, signupDate time.Time) []Marker {
	var markers []Marker

	for _, markerName := range markerNames {
		var marker Marker
		if err := db.Database.Where("name = ?", markerName).First(&marker).Error; err != nil {
			log.Printf("Error finding marker: %v", err)
			// Don't return an error if a marker is not found, just skip it
			continue
		}
		markers = append(markers, marker)
	}

	// Check if the user is created before July 1, 2023
	// If yes, append the Alpha marker to the markers slice
	alphaCutoffDate := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	if signupDate.Before(alphaCutoffDate) {
		var alphaMarker Marker
		if err := db.Database.Where("name = ?", "Alpha").First(&alphaMarker).Error; err != nil {
			log.Printf("Error finding Alpha marker: %v", err)
		} else {
			markers = append(markers, alphaMarker)
		}
	}

	return markers
}
