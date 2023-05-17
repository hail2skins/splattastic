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

// TeamTypeNew is a controller for creating a new team_type
func TeamTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"teamtypes/new.html",
		gin.H{
			"title":     "New Team Type",
			"header":    "New Team Type",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// TeamTypeCreate is a controller for creating a new team_type
func TeamTypeCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	models.TeamTypeCreate(name)
	c.Redirect(http.StatusFound, "/admin/teamtypes")
}

// TeamTypesIndex is a controller for the team_types index page
func TeamTypesIndex(c *gin.Context) {
	teamTypes, err := models.TeamTypesGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"teamtypes/index.html",
		gin.H{
			"title":     "Team Types",
			"header":    "Team Types",
			"teamTypes": teamTypes,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// TeamTypeShow is a controller for the team_types show page
func TeamTypeShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamType, err := models.TeamTypeShow(id)
	if err != nil {
		log.Printf("Error getting team_type: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"teamtypes/show.html",
		gin.H{
			"title":     "Team Type",
			"header":    "Team Type",
			"teamType":  teamType,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// TeamTypeEdit is a controller for the team_types edit page
func TeamTypeEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamType, err := models.TeamTypeShow(id)
	if err != nil {
		log.Printf("Error getting team_type: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"teamtypes/edit.html",
		gin.H{
			"title":     "Edit Team Type",
			"header":    "Edit Team Type",
			"teamType":  teamType,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}
