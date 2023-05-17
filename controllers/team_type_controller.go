package controllers

import (
	"net/http"
	"os"

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
