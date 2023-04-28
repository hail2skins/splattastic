package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
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

// EventTypeCreate is the controller for creating a new event type
func EventTypeCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.EventTypeCreate(name)
	c.Redirect(http.StatusFound, "/admin/eventtypes")
}

// EventTypesIndex is the controller for the event types index page
func EventTypesIndex(c *gin.Context) {
	eventTypes, err := models.EventTypesGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"eventtypes/index.html",
		gin.H{
			"title":      "Event Types",
			"header":     "Event Types",
			"logged_in":  h.IsUserLoggedIn(c),
			"eventtypes": eventTypes,
			"test_run":   os.Getenv("TEST_RUN") == "true",
		},
	)
}
