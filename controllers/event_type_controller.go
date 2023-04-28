package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"

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

// EventTypeShow is the controller for the event type show page
func EventTypeShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	eventType, err := models.EventTypeShow(id)
	if err != nil {
		log.Printf("Error getting event type: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"eventtypes/show.html",
		gin.H{
			"title":     "Event Type",
			"header":    "Event Type",
			"logged_in": h.IsUserLoggedIn(c),
			"eventtype": eventType,
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// EventTypeEdit is the controller for the event type edit page
func EventTypeEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	eventType, err := models.EventTypeShow(id)
	if err != nil {
		log.Printf("Error getting event type: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"eventtypes/edit.html",
		gin.H{
			"title":     "Edit Event Type",
			"header":    "Edit Event Type",
			"logged_in": h.IsUserLoggedIn(c),
			"eventtype": eventType,
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// EventTypeUpdate is the controller for updating an event type
func EventTypeUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	eventType, err := models.EventTypeShow(id)
	if err != nil {
		log.Printf("Error getting event type: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = eventType.Update(name)
	if err != nil {
		log.Printf("Error updating event type: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/eventtypes")
}
