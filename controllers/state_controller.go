package controllers

import (
	"net/http"
	"os"

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
