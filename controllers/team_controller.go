package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// TeamNew renders the new team page.
// Requires a team_type and a state
func TeamNew(c *gin.Context) {
	// Retrieve list of states for use in the form
	states, err := models.StatesGet()
	if err != nil {
		log.Printf("Error fetching states: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve list of team_types for use in the form
	teamTypes, err := models.TeamTypesGet()
	if err != nil {
		log.Printf("Error fetching team types: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"teams/new.html",
		gin.H{
			"title":      "New Team",
			"header":     "New Team",
			"logged_in":  h.IsUserLoggedIn(c),
			"states":     states,
			"team_types": teamTypes,
			"test_run":   os.Getenv("TEST_RUN") == "true",
			"user_id":    c.GetUint("user_id"),
		},
	)
}

// TeamCreate creates a new team
func TeamCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		log.Printf("Error creating team: name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	street := c.PostForm("street")
	street1 := c.PostForm("street1")
	city := c.PostForm("city")
	zip := c.PostForm("zip")
	if zip == "" {
		log.Printf("Error creating team: zip is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "zip is required"})
		return
	}

	abbreviation := c.PostForm("abbreviation")
	teamTypeIDStr := c.PostForm("team_type_id")
	teamTypeID, err := strconv.ParseUint(teamTypeIDStr, 10, 64)
	if err != nil {
		log.Printf("Error creating team: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stateIDStr := c.PostForm("state_id")
	stateID, err := strconv.ParseUint(stateIDStr, 10, 64)
	if err != nil {
		log.Printf("Error creating team: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from session claims
	userID := uint64(c.MustGet("user_id").(uint))

	// Check the associate_user checkbox and assign it to a variable
	associateUser := c.PostForm("associate_user")

	// Create the team
	team, err := models.TeamCreate(name, street, street1, city, zip, abbreviation, teamTypeID, stateID)
	if err != nil {
		log.Printf("Error creating team: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Associate the user with the team if the checkbox was checked
	if associateUser == "on" {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			log.Printf("Error associating user with team: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}
		userID := uint64(userIDVal.(uint))
		err = models.UserTeamCreate(uint64(team.ID), userID)
		if err != nil {
			log.Printf("Error associating user with team: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d/", userID))

}
