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
			"user_id":      c.GetUint("user_id"),
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

// DivesIndex renders the index page for dives
func DivesIndex(c *gin.Context) {
	dives, err := models.DivesGet()
	if err != nil {
		log.Printf("Error fetching dives: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"dives/index.html",
		gin.H{
			"title":     "Dives",
			"header":    "Dives",
			"logged_in": h.IsUserLoggedIn(c),
			"dives":     dives,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// DiveShow renders the show page for a dive
// Requires: DiveGet, DiveGroupGet, DiveTypeGet, BoardTypeGet, BoardHeightGet
func DiveShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	dive, err := models.DiveShow(id)
	if err != nil {
		log.Printf("Error fetching dive: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	diveGroup, _ := models.DiveGroupShow(dive.DiveGroupID)
	diveType, _ := models.DiveTypeShow(dive.DiveTypeID)
	boardType, _ := models.BoardTypeShow(dive.BoardTypeID)
	boardHeight, _ := models.BoardHeightShow(dive.BoardHeightID)

	c.HTML(
		http.StatusOK,
		"dives/show.html",
		gin.H{
			"title":       "Dive",
			"header":      "Dive",
			"logged_in":   h.IsUserLoggedIn(c),
			"dive":        dive,
			"divegroup":   diveGroup,
			"divetype":    diveType,
			"boardtype":   boardType,
			"boardheight": boardHeight,
			"test_run":    os.Getenv("TEST_RUN") == "true",
			"user_id":     c.GetUint("user_id"),
		},
	)
}

// DiveEdit renders the edit page for a dive
// Requires: DiveGet, DiveGroupGet, DiveTypeGet, BoardTypeGet, BoardHeightGet
func DiveEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	dive, err := models.DiveShow(id)
	if err != nil {
		log.Printf("Error fetching dive: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	diveGroup, _ := models.DiveGroupShow(dive.DiveGroupID)
	diveType, _ := models.DiveTypeShow(dive.DiveTypeID)
	boardType, _ := models.BoardTypeShow(dive.BoardTypeID)
	boardHeight, _ := models.BoardHeightShow(dive.BoardHeightID)

	diveGroups, _ := models.DiveGroupsGet()
	diveTypes, _ := models.DiveTypesGet()
	boardTypes, _ := models.BoardTypesGet()
	boardHeights, _ := models.GetBoardHeights()

	c.HTML(
		http.StatusOK,
		"dives/edit.html",
		gin.H{
			"title":        "Edit Dive",
			"header":       "Edit Dive",
			"logged_in":    h.IsUserLoggedIn(c),
			"dive":         dive,
			"divegroup":    diveGroup,
			"divetype":     diveType,
			"boardtype":    boardType,
			"boardheight":  boardHeight,
			"divegroups":   diveGroups,
			"divetypes":    diveTypes,
			"boardtypes":   boardTypes,
			"boardheights": boardHeights,
			"test_run":     os.Getenv("TEST_RUN") == "true",
			"user_id":      c.GetUint("user_id"),
		},
	)
}

// DiveUpdate updates a dive
func DiveUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		log.Printf("Error updating dive: name is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}
	numberStr := c.PostForm("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		log.Printf("Error converting number to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number"})
		return
	}
	difficultyStr := c.PostForm("difficulty")
	difficulty, err := strconv.ParseFloat(difficultyStr, 32)
	if err != nil {
		log.Printf("Error converting difficulty to float32: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid difficulty"})
		return
	}

	diveTypeIDStr := c.PostForm("divetype_id")
	diveTypeID, err := strconv.ParseUint(diveTypeIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting divetype_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid divetype_id"})
		return
	}

	diveGroupIDStr := c.PostForm("divegroup_id")
	diveGroupID, err := strconv.ParseUint(diveGroupIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting dive_group_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid divegroup_id"})
		return
	}
	boardTypeIDStr := c.PostForm("boardtype_id")
	boardTypeID, err := strconv.ParseUint(boardTypeIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting board_type_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boardtype_id"})
		return
	}
	boardHeightIDStr := c.PostForm("boardheight_id")
	boardHeightID, err := strconv.ParseUint(boardHeightIDStr, 10, 64)
	if err != nil {
		log.Printf("Error converting board_height_id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boardheight_id"})
		return
	}

	dive, err := models.DiveShow(id)
	if err != nil {
		log.Printf("Error fetching dive: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = dive.Update(name, number, float32(difficulty), diveTypeID, diveGroupID, boardTypeID, boardHeightID)
	if err != nil {
		log.Printf("Error updating dive: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/dives")

}

// DiveDelete deletes a dive
func DiveDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting id to uint64: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = models.DiveDelete(id)
	if err != nil {
		log.Printf("Error deleting dive: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/dives")
}
