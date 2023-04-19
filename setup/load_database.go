package setup

import (
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func LoadDatabase() {
	db.Connect()
	db.Database.AutoMigrate(&models.User{})
}
