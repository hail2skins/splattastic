package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/models"
)

func GetUserFromRequest(c *gin.Context) *models.User {
	userID := c.GetUint("user_id")

	var currentUser *models.User
	if userID > 0 {
		currentUser, _ = models.UserFind(uint64(userID))
	} else {
		currentUser = nil
	}
	return currentUser
}
