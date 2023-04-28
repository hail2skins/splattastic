package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

// EventTypeNew is the controller for the new event type page
func EventTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"eventtypes/new.html",
		gin.H{
			"title":     "New Event Type",
			"header":    "New Event Type",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}
