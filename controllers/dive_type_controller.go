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

// DiveTypeNew renders the new dive type form
func DiveTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"divetypes/new.html",
		gin.H{
			"title":     "New Dive Type",
			"header":    "New Dive Type",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// DiveTypeCreate creates a new dive type
func DiveTypeCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.DiveTypeCreate(name)
	c.Redirect(http.StatusFound, "/admin/divetypes")
}

// DiveTypesIndex gets all dive types
func DiveTypesIndex(c *gin.Context) {
	diveTypes, err := models.DiveTypesGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"divetypes/index.html",
		gin.H{
			"title":     "Dive Types",
			"header":    "Dive Types",
			"logged_in": h.IsUserLoggedIn(c),
			"divetypes": diveTypes,
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// DiveTypeShow renders the dive type show page
func DiveTypeShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	divetype, err := models.DiveTypeShow(id)
	if err != nil {
		log.Printf("Error getting dive type: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"divetypes/show.html",
		gin.H{
			"title":     "Dive Type",
			"header":    "Dive Type",
			"logged_in": h.IsUserLoggedIn(c),
			"divetype":  divetype,
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}
