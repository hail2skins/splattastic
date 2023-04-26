package controllers

import (
	"net/http"
	"os"

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
