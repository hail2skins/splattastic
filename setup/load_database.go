package setup

import (
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func LoadDatabase() {
	db.Connect()
	db.Database.AutoMigrate(&models.User{})
	db.Database.AutoMigrate(&models.UserType{})
	db.Database.AutoMigrate(&models.BoardHeight{})
	db.Database.AutoMigrate(&models.BoardType{})
	db.Database.AutoMigrate(&models.DiveType{})
	db.Database.AutoMigrate(&models.DiveGroup{})
	db.Database.AutoMigrate(&models.Dive{})
	db.Database.AutoMigrate(&models.EventType{})
	db.Database.AutoMigrate(&models.Event{})
	db.Database.AutoMigrate(&models.UserEventDive{})
	db.Database.AutoMigrate(&models.Score{})
	db.Database.AutoMigrate(&models.Marker{})
	db.Database.AutoMigrate(&models.UserMarker{})
	db.Database.AutoMigrate(&models.TeamType{})
}
