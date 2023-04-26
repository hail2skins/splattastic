package controllers

import (
	"net/http"
	"os"

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
		},
	)
}
