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

// DiveCreate handles the creation of a new dive
// Requires: DiveGroupGet, DiveTypeGet, BoardTypeGet, BoardHeightGet
func DiveCreate(c *gin.Context) {
	// Parse the form data
	name := c.PostForm("name")
	number, err := strconv.Atoi(c.PostForm("number"))
	if err != nil {
		log.Printf("Error converting number to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number"})
		return
	}
	difficulty, err := strconv.ParseFloat(c.PostForm("difficulty"), 32)
	if err != nil {
		log.Printf("Error converting difficulty to float32: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid difficulty"})
		return
	}
	divetypeID, err := strconv.ParseUint(c.PostForm("divetype_id"), 10, 64)
	if err != nil {
		log.Printf("Error converting divetype_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid divetype_id"})
		return
	}
	divegroupID, err := strconv.ParseUint(c.PostForm("divegroup_id"), 10, 64)
	if err != nil {
		log.Printf("Error converting divegroup_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid divegroup_id"})
		return
	}
	boardtypeID, err := strconv.ParseUint(c.PostForm("boardtype_id"), 10, 64)
	if err != nil {
		log.Printf("Error converting boardtype_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boardtype_id"})
		return
	}
	boardheightID, err := strconv.ParseUint(c.PostForm("boardheight_id"), 10, 64)
	if err != nil {
		log.Printf("Error converting boardheight_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boardheight_id"})
		return
	}

	// Create the new dive
	dive, err := models.DiveCreate(name, number, float32(difficulty), divetypeID, divegroupID, boardtypeID, boardheightID)
	if err != nil {
		log.Printf("Error creating dive: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirect to the newly created dive's show page
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/dives/%d", dive.ID))
}
