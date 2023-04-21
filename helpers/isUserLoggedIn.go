package helpers

import "github.com/gin-gonic/gin"

func IsUserLoggedIn(c *gin.Context) bool {
	userID, exists := c.Get("user_id")
	return exists && userID != nil
}
