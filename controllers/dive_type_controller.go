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
	letter := c.PostForm("letter")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.DiveTypeCreate(name, letter)
	c.Redirect(http.StatusFound, "/admin/divetypes")
}

// DiveTypesIndex gets all dive types
func DiveTypesIndex(c *gin.Context) {
	diveTypes, err := models.DiveTypesGet()
	if err != nil {
		log.Printf("Error fetching dive types: %v", err)
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

// DiveTypeEdit renders the dive type edit page
func DiveTypeEdit(c *gin.Context) {
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
		"divetypes/edit.html",
		gin.H{
			"title":     "Edit Dive Type",
			"header":    "Edit Dive Type",
			"logged_in": h.IsUserLoggedIn(c),
			"divetype":  divetype,
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// DiveTypeUpdate updates a dive type
func DiveTypeUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	letter := c.PostForm("letter")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	divetype, err := models.DiveTypeShow(id)
	if err != nil {
		log.Printf("Error getting dive type: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = divetype.Update(name, letter)
	if err != nil {
		log.Printf("Error updating dive type: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/divetypes")
}

// DiveTypeDelete deletes a dive type
func DiveTypeDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing dive type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = models.DiveTypeDelete(id)
	if err != nil {
		log.Printf("Error deleting dive type: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/divetypes")
}
