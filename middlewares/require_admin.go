package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/helpers"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is logged in
		if !helpers.IsUserLoggedIn(c) {
			// Redirect to login page with alert message
			c.Redirect(http.StatusSeeOther, "/?alert=You%20must%20be%20logged%20in%20to%20access%20this%20page")
			return
		}

		// Check if user is an admin
		isAdmin, exists := c.Get("isAdmin")
		if !exists {
			// Handle error
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !isAdmin.(bool) {
			// Redirect to home page with alert message
			c.Redirect(http.StatusSeeOther, "/?alert=You%20are%20forbidden%20from%20accessing%20this%20information")
			return
		}

		// Call the next handler
		c.Next()
	}
}
