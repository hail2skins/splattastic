package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
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
