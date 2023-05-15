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

// MarkerNew is a controller for creating a new marker
func MarkerNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"markers/new.html",
		gin.H{
			"title":     "New Marker",
			"header":    "New Marker",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// MarkerCreate is a controller for creating a new marker
func MarkerCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.MarkerCreate(name)
	c.Redirect(http.StatusFound, "/admin/markers")
}

// MarkersIndex is a controller for the markers index page
func MarkersIndex(c *gin.Context) {
	markers, err := models.MarkersGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"markers/index.html",
		gin.H{
			"title":     "Markers",
			"header":    "Markers",
			"logged_in": h.IsUserLoggedIn(c),
			"markers":   markers,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// MarkerShow is a controller for the marker show page
func MarkerShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	marker, err := models.MarkerShow(id)
	if err != nil {
		log.Printf("Error getting marker: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"markers/show.html",
		gin.H{
			"title":     "Marker",
			"header":    "Marker",
			"logged_in": h.IsUserLoggedIn(c),
			"marker":    marker,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// MarkerEdit is a controller for the marker edit page
func MarkerEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	marker, err := models.MarkerShow(id)
	if err != nil {
		log.Printf("Error getting marker: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"markers/edit.html",
		gin.H{
			"title":     "Edit Marker",
			"header":    "Edit Marker",
			"logged_in": h.IsUserLoggedIn(c),
			"marker":    marker,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// MarkerUpdate is a controller for updating a marker
func MarkerUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	marker, err := models.MarkerShow(id)
	if err != nil {
		log.Printf("Error getting marker: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = marker.Update(name)
	if err != nil {
		log.Printf("Error updating marker: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/markers")
}
