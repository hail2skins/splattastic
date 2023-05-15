package controllers

import (
	"net/http"
	"os"

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
