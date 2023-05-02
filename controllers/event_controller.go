package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ch "github.com/hail2skins/splattastic/controllers/helpers"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// EventNew renders the new page for entering an event tied to a user
func EventNew(c *gin.Context) {
	// Need to build the user relationship here as user is a FK in the event table
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ch.RenderEventErrorPage(c, "Error parsing user ID")
		return
	}
	// Retrieve the user from the database
	user, err := models.UserShow(id)
	if err != nil || user == nil {
		ch.RenderEventErrorPage(c, "Error retrieving user")
		return // Make sure to return after rendering the error page
	}

	// Build the event type dropdown here
	eventTypes, err := models.EventTypesGet()
	if err != nil {
		ch.RenderEventErrorPage(c, "Error retrieving event types")
	}

	// An event needs to select from dives, which can be added here
	// but can also be added from the event edit page during a meet entry
	dives, err := models.DivesGet()
	if err != nil {
		ch.RenderEventErrorPage(c, "Error retrieving dives")
	}

	c.HTML(
		http.StatusOK,
		"events/new.html", // Routes to /user/events/new
		gin.H{
			"title":      "New Event",
			"logged_in":  h.IsUserLoggedIn(c),
			"header":     "New Event",
			"eventTypes": eventTypes,
			"user_id":    c.GetUint("user_id"),
			"test_run":   os.Getenv("TEST_RUN") == "true",
			"dives":      dives,
			"user":       user,
		},
	)
}

// EventCreate handles the creation of a new event
// Requires Dives/Users/EventTypes
func EventCreate(c *gin.Context) {
	// Parse the form data
	name := c.PostForm("name")
	eventTypeIDStr := c.PostForm("event_type_id")
	eventTypeID, err := strconv.ParseUint(eventTypeIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event type ID to uint: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type ID"})
		return
	}
	against := c.PostForm("against")
	dateStr := c.PostForm("event_date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	location := c.PostForm("location")

	// Get the user ID from session claims
	userID := uint64(c.MustGet("user_id").(uint))

	// Handle dives
	diveIDsStr := c.PostForm("dive_ids")
	var diveIDs []uint64
	if diveIDsStr != "" {
		diveIDStrs := strings.Split(diveIDsStr, ",")
		diveIDs = make([]uint64, len(diveIDStrs))
		for i, diveIDStr := range diveIDStrs {
			diveID, err := strconv.ParseUint(diveIDStr, 10, 64)
			if err != nil {
				log.Printf("Error converting dive ID to uint: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dive ID"})
				return
			}
			diveIDs[i] = diveID
		}
	}

	// Create the event
	_, err = models.EventCreate(name, location, date, against, userID, eventTypeID, diveIDs)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating event"})
		return
	}

	// Redirect to the user show page, but will change this to the event show page
	// when that is complete
	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d", userID))
}
