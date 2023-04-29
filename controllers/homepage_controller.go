package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

func Home(c *gin.Context) {
	// Check if alert parameter exists in query string
	alert := c.Query("alert")
	if alert != "" {
		// Set the alert message to the template data
		c.Set("alert", alert)
	}

	c.HTML(
		http.StatusOK,
		"home/index.html",
		gin.H{
			"title":     "Splattastic",
			"logged_in": h.IsUserLoggedIn(c),
			"alert":     alert,
			"user_id":   c.GetUint("user_id"),
		})
}

func About(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home/about.html",
		gin.H{
			"title":     "About",
			"logged_in": h.IsUserLoggedIn(c),
			"user_id":   c.GetUint("user_id"),
		})
}
