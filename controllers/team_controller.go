package controllers

import (
	"log"
	"net/http"
	"os"

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
