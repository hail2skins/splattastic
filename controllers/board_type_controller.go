package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// BoardTypeNew renders the new board type form
func BoardTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"boardtypes/new.html",
		gin.H{
			"title":     "New Board Type",
			"logged_in": h.IsUserLoggedIn(c),
			"header":    "New Board Type",
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// BoardTypeCreate creates a new board type
func BoardTypeCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.BoardTypeCreate(name)
	c.Redirect(http.StatusFound, "/admin")
}
