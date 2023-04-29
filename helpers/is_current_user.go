package helpers

import "github.com/gin-gonic/gin"

// IsCurrentUser checks if the current user's ID matches the given user ID
func IsCurrentUser(c *gin.Context, userID uint64) bool {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		return false
	}
	return uint(userID) == sessionUserID.(uint)
}
