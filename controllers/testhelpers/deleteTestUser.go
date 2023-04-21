package testhelpers

import (
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func DeleteTestUser(email string) {
	user, err := models.GetUserByEmail(email)
	if err == nil && user != nil {
		db.Database.Unscoped().Delete(user)
	}
}
