package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	diveIDStrs := c.PostFormArray("dive_id")
	var diveIDs []uint64
	if len(diveIDStrs) > 0 {
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
	event, err := models.EventCreate(name, location, date, against, userID, eventTypeID, diveIDs)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating event"})
		return
	}

	// Redirect to the user show page, but will change this to the event show page
	// when that is complete
	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d/event/%d", userID, event.ID))
}

// EventShow renders the event show page
// Requires Users/EventTypes and UserEventDives
func EventShow(c *gin.Context) {
	// Get the event ID from the URL
	idStr := c.Param("event_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Retrieve the event from the database
	event, err := models.EventShow(id)
	if err != nil || event == nil {
		log.Printf("Error retrieving event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving event"})
		return
	}

	// Retrieve the user from the database
	user, err := models.UserShow(event.UserID)
	if err != nil || user == nil {
		log.Printf("Error retrieving user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user"})
		return
	}

	// Retrieve the event type from the database
	eventType, err := models.EventTypeShow(event.EventTypeID)
	if err != nil || eventType == nil {
		log.Printf("Error retrieving event type: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving event type"})
		return
	}

	// Retrieve the dives from the database
	dives, err := models.GetDivesForEvent(uint64(event.ID))
	if err != nil {
		log.Printf("Error retrieving dives for event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving dives"})
		return
	}

	// // Fetch the scores for each dive
	// diveScores := make(map[uint64][]models.Score)
	// for _, dive := range dives {
	// 	scores, err := models.FetchScores(uint64(user.ID), uint64(event.ID), uint64(dive.ID))
	// 	if err != nil {
	// 		log.Printf("Error fetching scores: %v", err)
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching scores"})
	// 		return
	// 	}
	// 	diveScores[uint64(dive.ID)] = scores
	// }

	// Retrieve the user event dives from the database
	userEventDives, err := models.GetUserEventDivesForEvent(id)
	if err != nil {
		log.Printf("Error retrieving user event dives for event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user event dives"})
		return
	}

	// Associate dives with user event dives
	diveIDToUserEventDiveID := make(map[uint64]uint64)
	for _, userEventDive := range userEventDives {
		diveIDToUserEventDiveID[userEventDive.DiveID] = uint64(userEventDive.ID)
	}

	// Format the event date
	formattedDate := event.Date.Format("01/02/2006")

	c.HTML(
		http.StatusOK,
		"events/show.html", // Routes to /user/events/:id
		gin.H{
			"title":                   "Event",
			"logged_in":               h.IsUserLoggedIn(c),
			"header":                  "Event",
			"event":                   event,
			"event_type":              eventType,
			"dives":                   dives,
			"diveIDToUserEventDiveID": diveIDToUserEventDiveID,
			"user":                    user,
			"user_id":                 c.GetUint("user_id"),
			"test_run":                os.Getenv("TEST_RUN") == "true",
			"current_user":            h.IsCurrentUser(c, uint64(user.ID)),
			"formatted_date":          formattedDate,
			//"scores":                  diveScores,
		},
	)
}

// GetUserEvents retrieves all user events for a given user
// Here is where we will pull in dives
func GetUserEvents(c *gin.Context) {
	// Get the user ID from the URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting user ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user for use in view
	user, err := models.UserShow(id)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user"})
		return
	}

	// Get the user events
	events, err := models.GetUserEvents(id)
	if err != nil {
		log.Printf("Error retrieving user events: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user events"})
		return
	}

	// Create a slice of maps to hold events and associated dives
	eventsWithDives := make([]map[string]interface{}, 0, len(events))

	// Loop through events to retrieve associated dives
	for _, event := range events {
		// Retrieve associated dives using model function
		dives, err := models.GetDivesForEvent(uint64(event.ID))
		if err != nil {
			log.Printf("Error retrieving dives for event %v: %v", event.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving dives for event"})
			return
		}

		eventWithDives := map[string]interface{}{
			"event": event,
			"dives": dives,
			"user":  user,
		}
		eventsWithDives = append(eventsWithDives, eventWithDives)
	}

	c.HTML(
		http.StatusOK,
		"events/user_index.html", // Routes to /user/:id/events
		gin.H{
			"title":        user.UserName + "'s Events",
			"logged_in":    h.IsUserLoggedIn(c),
			"header":       user.UserName + "'s Events",
			"events":       eventsWithDives,
			"test_run":     os.Getenv("TEST_RUN") == "true",
			"user_id":      c.GetUint("user_id"),
			"current_user": h.IsCurrentUser(c, uint64(user.ID)),
			"user":         user,
		},
	)
}

// EventEdit renders the event edit page
func EventEdit(c *gin.Context) {
	// Get the event ID from the URL
	idStr := c.Param("event_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Retrieve the event from the database
	event, err := models.EventShow(id)
	if err != nil || event == nil {
		log.Printf("Error retrieving event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving event"})
		return
	}

	// Retrieve the event type from the database
	eventType, err := models.EventTypeShow(event.EventTypeID)
	if err != nil || eventType == nil {
		log.Printf("Error retrieving event type: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving event type"})
		return
	}

	// Retrieve the dives from the database
	eventDives, err := models.GetDivesForEvent(uint64(event.ID))
	if err != nil {
		log.Printf("Error retrieving dives for event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving dives"})
		return
	}

	// retrieve event types from the database
	eventTypes, err := models.EventTypesGet()
	if err != nil {
		log.Printf("Error retrieving event types: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving event types"})
		return
	}

	// retrieve the dives from the database
	dives, err := models.DivesGet()
	if err != nil {
		log.Printf("Error retrieving dives: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving dives"})
		return
	}

	// Format the event date
	formattedDate := event.Date.Format("2006-01-02")

	c.HTML(
		http.StatusOK,
		"events/edit.html", // Routes to /user/events/:id/edit
		gin.H{
			"title":          "Edit Event",
			"logged_in":      h.IsUserLoggedIn(c),
			"header":         "Edit Event",
			"event":          event,
			"eventType":      eventType,
			"eventTypes":     eventTypes,
			"eventDives":     eventDives,
			"test_run":       os.Getenv("TEST_RUN") == "true",
			"user_id":        c.GetUint("user_id"),
			"current_user":   h.IsCurrentUser(c, uint64(event.UserID)),
			"formatted_date": formattedDate,
			"dives":          dives,
		},
	)
}

func EventUpdate(c *gin.Context) {
	// Get the event ID from the URL
	idStr := c.Param("event_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

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

	// Retrieve the event from the database
	event, err := models.EventShow(id)
	if err != nil || event == nil {
		log.Printf("Error retrieving event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving event"})
		return
	}

	// Handle dives
	diveIDStrs := c.PostFormArray("dive_id")
	var diveIDs []uint64
	if len(diveIDStrs) > 0 {
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

	// Update the event in the database
	_, err = event.Update(name, location, date, against, eventTypeID, diveIDs)
	if err != nil {
		log.Printf("Error updating event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating event"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d/event/%d", userID, id))
}

// Helper function to check if a uint64 slice contains a value
func contains(slice []uint64, value uint64) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// EventDelete deletes an event
func EventDelete(c *gin.Context) {
	// Get the event ID from the URL
	idStr := c.Param("event_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Get the user ID from session claims
	userID := uint64(c.MustGet("user_id").(uint))

	// Delete the event from the database
	err = models.EventDelete(id)
	if err != nil {
		log.Printf("Error deleting event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting event"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d", userID))
}

// EventScoreCreate creates a score for an event
// Can have between 1 and 9 judge scores
func EventScoreUpsert(c *gin.Context) {
	// Get the event ID from the URL
	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Get the user ID from session claims
	userID := uint64(c.MustGet("user_id").(uint))

	// Check if it's a POST request (form submission)
	if c.Request.Method == "POST" {
		var postData struct {
			UserID  uint64    `json:"userId"`
			EventID uint64    `json:"eventId"`
			DiveID  uint64    `json:"diveId"`
			Scores  []float64 `json:"scores"`
		}

		if err := c.BindJSON(&postData); err != nil {
			log.Printf("Error binding JSON data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Save the scores to the database
		for i, score := range postData.Scores {
			_, err := models.ScoreUpsert(postData.UserID, postData.EventID, postData.DiveID, i+1, score)
			if err != nil {
				log.Printf("Error creating score: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create score"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Scores saved successfully"})

		// Redirect to the event page
		c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d/event/%d", userID, eventID))
		return
	}

	c.HTML(
		http.StatusOK,
		"events/show.html",
		gin.H{
			"title":        "Create Score",
			"logged_in":    h.IsUserLoggedIn(c),
			"header":       "Create Score",
			"event_id":     eventID,
			"test_run":     os.Getenv("TEST_RUN") == "true",
			"user_id":      userID,
			"current_user": h.IsCurrentUser(c, userID),
		},
	)
}

// FetchScores fetches scores for an event for the JS
// and routes to /user/:user_id/event/:event_id/dive/:dive_id/scores
func FetchScores(c *gin.Context) {
	// Extract the IDs from the URL
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	eventID, _ := strconv.ParseUint(c.Param("event_id"), 10, 64)
	diveID, _ := strconv.ParseUint(c.Param("dive_id"), 10, 64)

	// Fetch the scores from the database
	scores, err := models.FetchScores(userID, eventID, diveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"scores": scores})
}

// EventDiveScore calculates the score for a dive in an event
func EventDiveScoreTotal(c *gin.Context) {
	// Get the dive ID from the URL
	diveIDStr := c.Param("dive_id")
	diveID, err := strconv.ParseUint(diveIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting dive ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dive ID"})
		return
	}

	// Total the scores
	total, err := models.CalculateDiveScore(diveID)
	if err != nil {
		log.Printf("Error calculating dive score: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating dive score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"score": total})
}

// EventMeetScore calculates the score for a meet, adding up all the dive scores
func EventMeetScore(c *gin.Context) {
	// Get the event ID from the URL
	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting event ID to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Total the meet score
	meetScore, err := models.CalculateMeetScore(eventID)
	if err != nil {
		log.Printf("Error calculating meet score: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating meet score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"score": meetScore})

}
