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

// StateNew handles the request to display the new state form
func StateNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"states/new.html",
		gin.H{
			"title":     "New State",
			"header":    "New State",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// StateCreate handles the request to create a new state
func StateCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	code := c.PostForm("code")
	if code == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	models.StateCreate(name, code)
	c.Redirect(http.StatusFound, "/admin/states")
}

// StatesIndex handles the request to display the states index page
func StatesIndex(c *gin.Context) {
	states, err := models.StatesGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"states/index.html",
		gin.H{
			"title":     "States",
			"header":    "States",
			"states":    states,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// StateShow handles the request to display a state
func StateShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	state, err := models.StateShow(id)
	if err != nil {
		log.Printf("Error getting state: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"states/show.html",
		gin.H{
			"title":     "State",
			"header":    "State",
			"state":     state,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// StateEdit handles the request to display the edit state form
func StateEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	state, err := models.StateShow(id)
	if err != nil {
		log.Printf("Error getting state: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(
		http.StatusOK,
		"states/edit.html",
		gin.H{
			"title":     "Edit State",
			"header":    "Edit State",
			"state":     state,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}
