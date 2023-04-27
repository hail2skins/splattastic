package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/models"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the session and sessionID
		session := sessions.Default(c)
		sessionID := session.Get("id")

		// Declare user and userPresent variables
		var user *models.User
		userPresent := true

		// Check if there is no sessionID
		if sessionID == nil {
			userPresent = false
		} else {
			// Find user by sessionID
			user, _ = models.UserFind(sessionID.(uint64))

			// Check if user was found
			userPresent = (user.ID > 0)
		}

		// Check if user is an admin
		isAdmin := user != nil && user.Admin != nil && *user.Admin

		// If user was found, set user_id, email, and isAdmin in the gin.Context
		if userPresent {
			c.Set("user_id", user.ID)
			c.Set("email", user.Email)
			c.Set("isAdmin", isAdmin)
		}

		// Set the logged_in value in the gin.Context for use in templates and controllers
		c.Set("logged_in", userPresent)

		// Call the next handler
		c.Next()
	}
}
