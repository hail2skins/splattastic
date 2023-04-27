package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// DiveNew renders the new dive form
// Requires: DiveGroupGet, DiveTypeGet, BoardTypeGet, BoardHeightGet
func DiveNew(c *gin.Context) {
	diveGroups, err := models.DiveGroupsGet()
	if err != nil {
		log.Printf("Error fetching dive groups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	diveTypes, err := models.DiveTypesGet()
	if err != nil {
		log.Printf("Error fetching dive types: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	boardTypes, err := models.BoardTypesGet()
	if err != nil {
		log.Printf("Error fetching board types: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	boardHeights, err := models.GetBoardHeights()
	if err != nil {
		log.Printf("Error fetching board heights: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"dives/new.html",
		gin.H{
			"title":        "New Dive",
			"header":       "New Dive",
			"logged_in":    h.IsUserLoggedIn(c),
			"divegroups":   diveGroups,
			"divetypes":    diveTypes,
			"boardtypes":   boardTypes,
			"boardheights": boardHeights,
			"test_run":     os.Getenv("TEST_RUN") == "true",
		},
	)
}
