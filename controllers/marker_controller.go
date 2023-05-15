package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
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
