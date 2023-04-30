package controllers

import (
	"net/http"
	"os"
	"strconv"

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
