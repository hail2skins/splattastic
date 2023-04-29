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

// DiveGroupNew renders the new dive group form
func DiveGroupNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"divegroups/new.html",
		gin.H{
			"title":     "New Dive Group",
			"header":    "New Dive Group",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// DiveGroupCreate creates a new dive group
func DiveGroupCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.DiveGroupCreate(name)
	c.Redirect(http.StatusFound, "/admin/divegroups")
}

// DiveGroupsIndex gets all dive groups
func DiveGroupsIndex(c *gin.Context) {
	diveGroups, err := models.DiveGroupsGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"divegroups/index.html",
		gin.H{
			"title":      "Dive Groups",
			"header":     "Dive Groups",
			"logged_in":  h.IsUserLoggedIn(c),
			"divegroups": diveGroups,
			"test_run":   os.Getenv("TEST_RUN") == "true",
			"user_id":    c.GetUint("user_id"),
		},
	)
}

// DiveGroupShow renders the dive group show page
func DiveGroupShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive group id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	diveGroup, err := models.DiveGroupShow(id)
	if err != nil {
		log.Printf("Error getting dive group: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"divegroups/show.html",
		gin.H{
			"title":     "Dive Group",
			"header":    "Dive Group",
			"logged_in": h.IsUserLoggedIn(c),
			"divegroup": diveGroup,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// DiveGroupEdit renders the dive group edit form
func DiveGroupEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive group id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	diveGroup, err := models.DiveGroupShow(id)
	if err != nil {
		log.Printf("Error getting dive group: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"divegroups/edit.html",
		gin.H{
			"title":     "Edit Dive Group",
			"header":    "Edit Dive Group",
			"logged_in": h.IsUserLoggedIn(c),
			"divegroup": diveGroup,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// DiveGroupUpdate updates a dive group
func DiveGroupUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive group id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	diveGroup, err := models.DiveGroupShow(id)
	if err != nil {
		log.Printf("Error getting dive group: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = diveGroup.Update(name)
	if err != nil {
		log.Printf("Error updating dive group: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/divegroups")
}

// DiveGroupDelete deletes a dive group
func DiveGroupDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive group id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = models.DiveGroupDelete(id)
	if err != nil {
		log.Printf("Error deleting dive group: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/divegroups")
}
