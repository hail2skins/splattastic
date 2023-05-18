package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
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
