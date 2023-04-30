package setup

import (
	"log"

	"github.com/hail2skins/splattastic/database"
)

// RemoveUniqueConstraintOnDiveNumber removes the unique constraint on the dive number
// Run this function once in a startup, then remove it from main
// But we will keep this here for future needs.
func RemoveUniqueConstraintOnDiveNumber() {
	db, err := database.Database.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	_, err = db.Exec("ALTER TABLE dives DROP CONSTRAINT dives_number_key")
	if err != nil {
		log.Fatalf("Failed to remove unique constraint on Dive number: %v", err)
	}
}
